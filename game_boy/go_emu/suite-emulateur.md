# Suite de l’émulateur

État : les 11 ROMs Blargg `cpu_instrs` passent. Tu as CPU, timer, IE/IF/IME, HALT, série.
Il manque : l’image, le son, les cartouches > 32 Ko, et une boucle qui fait tourner tout ça ensemble.

Cible : réécriture en Rust. L’émulateur produit **un tampon 160×144** et **des samples**.
La fenêtre et la carte son sont de la plomberie côté host, pas de l’émulation.

Ordre : boucle → PPU → joypad → MBC → APU.

Les pages citées sont celles de `GBCPUman.pdf` (numéro imprimé = page du PDF, 139 pages).
Ce manuel décrit surtout les **registres**. Deux sujets y sont incomplets, ne les cherche pas :
le timing du PPU au cycle près et le frame sequencer de l’APU (là, Pan Docs).

---

## Les chiffres à connaître

Tout le reste en découle.

- Horloge : 4 194 304 tics/s.
- Une scanline = **456 tics**.
- Une frame = **154 scanlines** = **70 224 tics** ≈ 59,7 images/s.
- Lignes 0–143 : visibles. Lignes 144–153 : VBlank.

**Manuel p. 54** : le seul endroit où ces valeurs sont écrites (456 clks par cycle de modes,
4560 pour le VBlank, 70224 pour un écran complet). C’est dans la description de STAT, pas dans §2.8.

Ton `UpdateTimer(n)` est déjà le bon modèle : « n tics viennent de passer ».
PPU et APU sont deux abonnés de plus au même signal.

---

## Étape 0 — La boucle

**Manuel p. 17–18** (valeurs des registres au démarrage : AF, SP, LCDC=$91, BGP=$FC, tous les NR),
**p. 34** (le VBlank arrive ~59,7 fois par seconde), **p. 54** (70224 clks par écran).
Initialise avec les valeurs de la p. 17–18 plutôt que tout à zéro : certains jeux les supposent.

Aujourd’hui le CPU exécute puis s’arrête. Il faut qu’il tourne sans fin, à la bonne vitesse.

Une itération = exécuter des instructions jusqu’à cumuler 70 224 tics, puis rendre la main.
Chaque instruction rend son nombre de tics ; tu les passes au timer, au PPU, à l’APU.
Le host attend le temps restant pour tenir ~60 Hz.

Sans écran, un mode headless suffit : N frames, puis tu dumpes l’image en PPM.

Preuve : les tests Blargg passent toujours, et une ROM tourne indéfiniment sans erreur.

---

## Étape 1 — Le compteur de lignes (PPU sans pixels)

**Manuel p. 55** (LY : valeurs 0 à 153, 144–153 = VBlank), **p. 34** (l’interruption VBlank),
**p. 39** (IF, bit 0 = VBlank), **p. 40** (table des vecteurs : VBlank → $0040).

Le plus important, et ça ne dessine rien.

Tu accumules les tics. À 456, tu retires 456 et tu incrémentes **LY** (`$FF44`).
Quand LY atteint 144 : tu lèves **IF bit 0** (VBlank, vecteur `$0040`).
Quand LY atteint 154 : LY repart à 0, la frame est finie.

Pourquoi d’abord : quasi tous les jeux bouclent sur « attendre LY == 144 » avant de toucher la VRAM.
Sans ce compteur, ils sont figés pour toujours et tu crois que ton CPU est cassé.

Preuve : une ROM qui affichait un écran noir se met à exécuter du code neuf ; LY défile.

---

## Étape 2 — Décoder une tuile

**Manuel p. 24** : le schéma d’une tuile octet par octet, avec le résultat en pixels.
**p. 57** pour BGP.

Une tuile = 8×8 pixels, **16 octets**, soit 2 octets par ligne.
Ces 2 octets sont des **plans de bits** : le premier donne le bit de poids faible de chaque pixel,
le second le bit de poids fort. Pixel = valeur 0 à 3.

Ces 0–3 ne sont pas des couleurs : ce sont des index dans la palette **BGP** (`$FF47`),
qui associe chacun à un des 4 gris.

Preuve : tu extrais une tuile de la VRAM et tu retrouves une forme lisible (une lettre du logo).

---

## Étape 3 — Le fond, ligne par ligne

**Manuel p. 22–23** (fond de 256×256, tilemap 32×32, l’enroulement, les deux tables de tuiles),
**p. 51–52** (LCDC bit par bit), **p. 54** (SCY), **p. 55** (SCX), **p. 57** (BGP),
**p. 8–9** (plan mémoire, où est la VRAM).

Pour la ligne LY, tu calcules 160 pixels.

La tilemap (32×32 index de tuiles) est en `$9800` ou `$9C00` selon **LCDC bit 3** (`$FF40`).
La position dans le fond est `LY + SCY` verticalement, `x + SCX` horizontalement, **modulo 256** : le fond boucle.
Diviser par 8 donne la case de la tilemap, le reste donne le pixel dans la tuile.

Piège classique : **LCDC bit 4** choisit l’adressage des tuiles.
À 1, base `$8000` et index non signé. À 0, base `$9000` et index **signé** (−128 à 127).
Se tromper ici donne des graphismes corrects mais mélangés.

Si **LCDC bit 7** est à 0, l’écran est éteint : tampon blanc, LY figé à 0.

Preuve : le logo Nintendo scrolle tout seul, ou le fond de Tetris s’affiche.
Le dessin manuel du logo dans `Boot()` devient inutile — supprime-le.

---

## Étape 4 — Fenêtre et sprites

**Fenêtre** : p. 23 (comportement) et **p. 58–59** (WY, WX, avec un schéma de placement).
**Sprites** : **p. 25–27** (les 40 blocs de 4 octets, les 10 par ligne, les décalages −8/−16,
les priorités, le détail des flags du 4e octet), **p. 58** (OBP0/OBP1).
**STAT** : **p. 52–54**. **DMA OAM** ($FF46) : **p. 55–57**, presque tous les jeux l’utilisent
pour remplir l’OAM pendant le VBlank.

La **fenêtre** est une seconde couche opaque, non scrollée, qui démarre en WX−7 / WY.
Utilisée pour les barres de score. Elle a son propre compteur de lignes interne.

Les **sprites** sont dans l’OAM (`$FE00`–`$FE9F`), 40 entrées de 4 octets : Y, X, index de tuile, attributs.
Règles : Y est décalé de 16 et X de 8 ; **10 sprites maximum par ligne** ;
la couleur 0 est transparente ; un bit d’attribut met le sprite derrière le fond.

Preuve : `dmg-acid2` affiche la tête, ou Mario apparaît sur son décor.

Le registre **STAT** (`$FF41`) et son interruption servent aux effets par ligne (split screen).
LYC == LY suffit au début ; les modes 0/2/3 exacts, plus tard.

---

## Étape 5 — Joypad

**Manuel p. 35–37** : le registre P1, le schéma de la matrice P14/P15, et un exemple réel
tiré de Ms. Pacman qui montre exactement la séquence attendue par un jeu.
**p. 35** pour l’interruption (front haut→bas de P10–P13), **p. 39** (IF bit 4).

Un seul registre, `$FF00`, et il est **inversé** : 0 = pressé.

Les bits 4 et 5 sélectionnent quel groupe tu lis : direction ou boutons.
Le jeu écrit un sélecteur, puis relit les 4 bits bas.
Tu réponds avec l’état des touches du groupe sélectionné.

Preuve : « PRESS START » réagit.

---

## Étape 6 — Cartouches (MBC)

**Manuel p. 11** (la liste des types à $0147), **p. 12** (tailles ROM à $0148 et RAM à $0149),
**p. 13–14** (MBC1 en entier : les quatre plages d’écriture et le mode 16/8 contre 4/32),
**p. 14–15** (MBC2), **p. 15** (MBC3), **p. 15–16** (MBC5).
Le plan mémoire général est **p. 8–9**.

Blargg tient dans 32 Ko. Presque aucun jeu commercial.

L’en-tête à `$0147` donne le type de mapper, `$0148` la taille ROM, `$0149` la taille RAM.
Commence par lire ça et l’afficher : tu sauras tout de suite si une ROM est jouable chez toi.

Le principe du MBC : `$4000`–`$7FFF` n’est pas fixe, c’est une **fenêtre** sur une banque de 16 Ko.
Écrire en ROM ne modifie pas la ROM, ça **pilote le mapper**.

MBC1, les quatre zones d’écriture :
- `$0000`–`$1FFF` : valeur `$0A` active la RAM de sauvegarde en `$A000`, `$00` la coupe.
- `$2000`–`$3FFF` : numéro de banque ROM (5 bits ; les valeurs 0 et 1 donnent toutes deux la banque 1).
- `$4000`–`$5FFF` : 2 bits de plus, dont le sens dépend du mode ci-dessous.
- `$6000`–`$7FFF` : le mode. 0 = beaucoup de ROM / 8 Ko de RAM, 1 = moins de ROM / 32 Ko de RAM.
  En mode 0 ces 2 bits sont les bits hauts de la banque ROM, en mode 1 ils choisissent la banque RAM.

C’est le moment de supprimer l’erreur « MBC non implémenté » de ton écriture mémoire.

Preuve : Super Mario Land démarre. Tetris marche déjà (ROM only).

Sauvegarde `.sav` : simple dump de la RAM cartouche, une fois qu’elle fonctionne.

---

## Étape 7 — Son (APU)

**Manuel p. 28–29** : vue d’ensemble des 4 voies, et surtout la conversion Hz ↔ registre
(`gb = 2048 − 131072/Hz`). Puis un registre par page :
NR10 sweep **p. 40–41**, NR11 duty/length **p. 41**, NR12 envelope **p. 42**, NR13/NR14 **p. 42–43**,
voie 2 (NR21–NR24) **p. 43–45**, voie 3 (NR30–NR34) **p. 45–47**,
voie 4 (NR41–NR44) **p. 47–49**, NR50 et NR51 **p. 50**, NR52 **p. 51**, wave RAM **p. 51**.

Ce que le manuel **ne dit pas** : le frame sequencer à 512 Hz. Il donne seulement les résultats
(sweep en 1/128 s p. 40, envelope en 1/64 s p. 42, length en 1/256 s p. 41). C’est pour ça
que le calage précis se lit ailleurs.

Tu as les registres NR10–NR52 et la wave RAM. Il manque l’horloge et la sortie.

Deux horloges, ne les confonds pas :

1. **Le frame sequencer**, tous les 8 192 tics (512 Hz). Il fait vieillir les notes :
   la durée (256 Hz), le volume/envelope (64 Hz), le sweep (128 Hz).
2. **Le timer de fréquence** de chaque voie, qui fait la hauteur du son.
   Sa période vaut (2048 − fréquence) × 4 tics ; à chaque échéance tu avances d’un pas
   dans le motif de 8 valeurs du duty cycle.

Écrire dans le registre haut avec le bit 7 = **trigger** : la note (re)démarre.

Sortie : tu prends un échantillon tous les ~95 tics pour du 44,1 kHz.
NR51 décide gauche/droite, NR50 le volume master.

Commence par **la voie 2 seule** (carré, sans sweep) : c’est la plus simple et elle valide
toute la chaîne tics → APU → haut-parleur. Puis voie 1 (sweep), voie 4 (bruit), voie 3 (wave).

Preuve : un bip audible. Les tests `dmg_sound` de Blargg viennent après, pas avant.

---

## À ignorer pour l’instant

GBC (double vitesse, couleurs), SGB, timing PPU au cycle près,
câble série (p. 31, 37 — déjà suffisant pour Blargg), bug OAM (p. 27),
mode STOP (p. 19).

Quand un jeu ne démarre pas, cherche dans cet ordre :
banque ROM manquante → LY/VBlank → joypad → le son.
