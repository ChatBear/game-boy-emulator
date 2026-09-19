# Cycles Game Boy (GBCPUman)

Le Game Boy ne fait pas « une instruction, puis une frame ». Tout partage **une seule horloge**. Le cycle est cette unité de temps commune.

## Deux noms, la même horloge

CPU à **4,194304 MHz**. GBCPUman compte en **cycles d’horloge**. Nintendo compte en **cycles machine** :

| | Fréquence | `NOP` |
|---|---|---|
| Cycle machine | ~1,05 MHz | 1 |
| Cycle d’horloge | ~4,19 MHz | 4 |

**1 cycle machine = 4 cycles d’horloge.** Chaque coût d’instruction dans le manuel est un multiple de 4.

Sur le bus (annexe 5.2), un cycle machine = un accès mémoire (adresse, puis lecture ou écriture). D’où un CPU « à ~1 MHz » pour le travail, même si l’oscillateur est à 4,19 MHz.

## Une étape CPU

1. **Fetch** de l’opcode à `PC` (c’est aussi là que les interruptions sont testées).
2. **Decode / execute** — octets supplémentaires si besoin (`n`, `nn`, préfixe `CB`).
3. **Avancer le temps** du nombre de cycles de l’instruction.

Exemples :

- `NOP`, `LD A,B`, `INC A` → **4** (1 accès : l’opcode)
- `LD A,n`, `LD A,(HL)` → **8** (opcode + une lecture)
- `JP nn` dans ce manuel → **12**
- `LD (nn),SP` → **20**

Règle : chaque accès mémoire extra coûte 4 cycles. Ops registres = pas cher ; `(HL)`, immédiat, pile = plus cher.

`PC` démarre à `$0100` après boot. Ensuite seules les instructions le bougent : incrément séquentiel, ou `JP` / `JR` / `CALL` / `RST` / interruption.

## Le même compteur fait tourner le reste

Après (ou pendant) une instruction, **tout le reste avance du même nombre de cycles d’horloge**.

**PPU (LCD), registre STAT `$FF41` :**

- Une ligne = **456** cycles
- Modes 2 → 3 → 0 sur une ligne visible (~80 + ~172 + ~204 ; plages 77–83 / 169–175 / 201–207)
- 144 lignes visibles, puis VBlank (mode 1) : **10 lignes = 4560** cycles
- Une frame = **70224** cycles ≈ **59,73 Hz**
- `LY` (`$FF44`) = ligne courante, 0–153 (144–153 = VBlank)

Les jeux qui bouclent sur `LY` ou STAT attendent un point précis de cette ligne de 456 cycles.

**DIV / timer :**

- `DIV` (`$FF04`) : **16384** fois/s → tous les **256** cycles. Toute écriture le remet à 0.
- `TIMA` : 4096 / 16384 / 65536 / 262144 Hz selon TAC, soit tous les **1024 / 256 / 64 / 16** cycles. Overflow → charge `TMA`, pose le bit IF timer.

**Série :** horloge interne à 8192 Hz (~122 µs/bit, 8 bits puis interruption).

**DMA :** écrire `$FF46` copie 160 octets vers l’OAM (~160 µs). Pendant le DMA, seule la HRAM est utilisable. La boucle d’attente du manuel (`dec a` = 4, `jr nz` = 12) sert à attendre la fin du transfert.

## Interruptions (fetch)

§2.12.1 :

1. Le matériel pose un bit dans `IF`.
2. Si `IME` **et** le bit `IE` correspondant sont set, au **prochain fetch** :
   - clear `IME`
   - push `PC`
   - jump `$0040` / `$0048` / `$0050` / `$0058` / `$0060`

Priorité si plusieurs en même temps : VBlank > LCD STAT > Timer > Série > Joypad.

`EI` / `DI` prennent effet **après l’instruction suivante**, pas tout de suite.

`HALT` arrête l’horloge CPU jusqu’à une interruption (le reste peut continuer). `STOP` arrête CPU et LCD jusqu’à un bouton.

## Boucle mentale de l’émulateur

Pas « CPU à fond, puis dessiner ». C’est :

```
fetch opcode
si interruption pendante (IME + IE + IF) → la servir à la place
exécuter l’instruction  → coûte N cycles d’horloge
avancer PPU de N
avancer timer / DIV / série / APU de N
répéter
```

`N` est toujours un multiple de 4. Synchroniser le reste tous les 4 cycles (1 cycle machine) suffit pour un premier émulateur ; un PPU pixel-accurate demandera plus fin.

## Limite de ce PDF

GBCPUman donne des cycles d’**horloge**, mais aplatit plusieurs timings **conditionnels** (`JR cc` toujours 8, `JP cc` toujours 12, `CALL nn` 12). Sur le vrai hardware, un saut/call **pris** coûte des cycles machine en plus (accès mémoire extra). Les tables = ordre de grandeur ; vérifier une table d’opcodes plus récente pour les coûts pris / non pris.

Ce temps extra doit aussi aller au PPU et au timer. Un jeu qui attend HBlank ou utilise le timer se désynchronise si les coûts d’instructions sont faux.
