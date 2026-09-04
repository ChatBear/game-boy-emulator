# Problèmes de design du CPU

État : les 11 ROM Blargg `cpu_instrs` passent. Ce fichier ne parle pas de ce qui manque (PPU, son, cartouches). Il parle de **comment le code est construit**.

Le problème n’est pas « 11 fichiers ». C’est que le jeu d’instructions est traité comme **~500 cas particuliers**, classés comme les chapitres du GBCPUman, au lieu d’un **décodeur + quelques opérations**.

---

## 1. La taxonomie vient du manuel, pas de la puce

Les fichiers suivent le sommaire : load, alu, jumps, calls, returns, restarts, bit, rotate_shift, miscellaneous, stack, 16_bit_arithmetic.

Sur le CPU, les mêmes 3 bits servent partout :

| bits 0–2 | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 |
|---|---|---|---|---|---|---|---|---|
| registre | B | C | D | E | H | L | (HL) | A |

Ça encode `LD r,r′` (0x40–0x7F sauf HALT 0x76), l’ALU sur A (0x80–0xBF), toute la page CB, et (avec 2 bits) les conditions NZ/Z/NC/C.

`cpu_bit.go` a déjà compris ça (une boucle). Les autres fichiers écrivent la même boucle à la main. SWAP (CB 0x30–0x37) est dans miscellaneous alors que c’est le même motif que RLC.

Ouvrir un fichier ne répond pas à « que fait l’octet 0x80 ? ». Il répond à « chapitre ALU du PDF ». L’encodage n’existe nulle part dans le source.

Les 11 fichiers se videront **après** un vrai décodeur, pas parce qu’on les fusionne.

---

## 2. On ne réécrit pas le même opcode : on réécrit la même opération

Aucun numéro n’est (en principe) enregistré deux fois. Ce qui est recopié, c’est **l’opération + un registre**.

- `addA` existe une fois. Puis 9 closures : A, B, C, D, E, H, L, (HL), immédiat. Idem ADC, SUB, SBC, AND, OR, XOR, CP.
- `LD r,r′` : une grille 8×8. `LD B,B` ne fait que `UpdateTimer(4)`.
- RLC, RL, RRC, RR, SLA, SRA, SRL : chacun recopié 8 fois.
- SWAP : recopié 8 fois, alors que BIT/SET/RES à côté sont en boucle.
- NZ / Z / NC / C : recopié ~16 fois (JP, JR, CALL, RET).
- Pousser PC : `call`, `restart`, `serviceInterrupt`.
- `ADD SP,e` (0xE8) et `LD HL,SP+e` (0xF8) : le même calcul H/C, un dans `cpu_16_bit_arithmetic.go`, un dans `cpu_load.go`.

Deux opcodes différents peuvent être la même fonction + un argument. Le code ne le dit pas.

---

## 3. Les closures cachent le décodeur

`NewCPU` → `initOpcodes()` remplit deux tableaux `[256]func()`. Chaque entrée capture `cpu`.

- Pas de donnée d’instruction : ni nom, ni longueur, ni cycles. Les cycles sont un `UpdateTimer(n)` au fond de la closure.
- Plusieurs `n` sont faux (`JP` 12, `JR` 8, `RET` 8, `RST` 32, `LD A,B` 8). Blargg `cpu_instrs` ne le voit pas ; `instr_timing` le verrait.
- Cycles pris / pas pris : JP cc, JR cc, CALL cc, RET cc n’ont pas le même coût selon la condition. Tout le monde paie le même `n`.
- Une instance = 512 closures sur le tas. Ce n’est pas une ROM de décodage.
- Collision silencieuse : rien n’empêche deux `init*` d’écrire la même case. HALT 0x76 est dans miscellaneous parce que load a un trou. Accord humain, pas une règle.
- CB est un `if` dans `Step`. EI (`pendingEnableIME`) s’applique après un opcode normal et est sauté après un CB. Le délai devient « prochaine instruction non-CB ».

---

## 4. Pas d’API registres, donc pas de facteur commun

Il n’existe pas `readR(r)` / `writeR(r, v)`. Chaque handler nomme `cpu.b`, `cpu.c` à la main. Sans ça on ne peut pas écrire « ALU[op](readR(r)) » ni « writeR(dst, readR(src)) ».

`(HL)` n’est qu’un `r == 6` plus cher. Aujourd’hui c’est un cas spécial recopié, avec deux styles : `getHL()` vs `uint16(h)<<8|l`.

INC/DEC rr : `setBC(getBC()+1)` recopié au lieu d’un index de paire.

---

## 5. Les flags sont des hex magiques, pas un contrat

Partout : `0x80` (Z), `0x40` (N), `0x20` (H), `0x10` (C). Pas de constante, pas de `setZ`.

Trois politiques se mélangent : `addA` remet F à 0 ; `inc`/`bit` gardent C ; `add16HL` efface N/H/C et garde Z ; DAA garde N.

RLCA vs CB RLC : Z toujours 0 vs Z selon le résultat — **ça**, c’est le hardware. Le reste est du bricolage recopié.

API interne incohérente : `addA` mute A et F ; `inc`/`rlc` prennent une valeur, renvoient une valeur, mutent F en coulisse.

---

## 6. Le type `CPU` est un fourre-tout

Ce qui est le CPU : registres, PC, SP, IME, HALT, `Step`.

Ce qui n’en est pas : `memory` 64 Ko, `Screen`, `Palette`, `scx`/`scy` (morts), `apu` (jamais branché), `serialOutput`, timer (`cycle`/`timerAcc`), `Boot()` (~50 lignes de logo Nintendo).

`NewCPU(a,b,c,d,e,f,h,l)` : les tests passent des zéros. `InitializeRegisterValues` ne pose que PC et SP. Le constructeur ment sur ce qu’il initialise.

`tick` est un alias de `UpdateTimer`. Les handlers appellent l’un, l’IRQ l’autre.

---

## 7. `Step` mélange trop de responsabilités

Fetch, IRQ, HALT, EI et timer sont dans la même fonction que le dispatch. On ne peut pas tester « juste 0x80 ».

Les cycles ne remontent pas : chaque handler appelle le timer. `Step` ne renvoie pas `n`. La future boucle 70 224 tics n’a rien à écouter.

Commentaire de DI faux (le code coupe IME tout de suite ; le commentaire décrit EI). STOP ne consomme pas l’octet 0x00 qui le suit. `Game.Update` est vide : la fenêtre ne fait jamais tourner le CPU.

---

## 8. Deux chemins pour le même geste

Lire `(HL)` : parfois `memory[getHL()]`, parfois HL reconstruit à la main. Écrire : `writeMemory` (retourne `error`, tous font `_ =`) ou `memory[addr] =` (BIT SET/RES, certains rotates). Pousser 16 bits : `call`/`restart` vs PUSH recopié 4×. Flags SP+e : `add16SP` et copie dans 0xF8.

Une `error` n’est pas un modèle pour un bus : le hardware ignore ou a un effet de bord, il n’« échoue » pas.

---

## 9. La preuve, c’est une ROM — pas le code

`cpu_test.go` lance les 11 ROM, lit la série, cherche `"Passed"`. Bon test d’intégration. Seul test.

On ne peut pas écrire `ADD A,B` avec A=0x0F, B=0x01 → flags attendus. Les helpers ALU pourraient être testés seuls. Ils ne le sont pas.

Blargg `cpu_instrs` ne valide pas les timings, HALT bug, DIV, PPU, MBC. « Passed » dit que SB/SC coincide, pas que le décodeur est bien conçu.

---

## 10. Pas de langage commun dans le paquet

`init16BitArthmeticOpCode` vs `initJumpsOpCodes` vs `initBitOpCode`. `adress`. Commentaires EI/DI vagues. `nCycles` est une map pour 4 constantes. `Run` existe ; `main` ne l’appelle pas.

Sans vocabulaire stable (r, cc, flags, tics), chaque fichier réinvente le même if.

---

## À garder

`addA` et la famille ALU, `rlc`/`sla`/…, `bit`/`set`/`res`, `getAF`/`setAF`, `fetch8`/`fetch16`, IME/HALT/IRQ dans `Step` (le quoi, pas le mélange), la boucle de `cpu_bit.go`.

---

## Direction

1. `readR` / `writeR` et `cond(cc)`. Ça tue la grille LD, l’ALU recopié, CB, et les 16 if de condition.
2. Décoder dans `Step` (ou remplir les tables par boucle sur les bits). Les helpers restent.
3. `Step` renvoie les tics. Un seul endroit tick le timer (puis le PPU).
4. Flags nommés.
5. `CPU` = registres + IME + HALT. Logo, écran, palette, scx, APU sortent.
6. Un test par helper ALU en plus de Blargg.

Ne pas fusionner les 11 fichiers d’abord. Enlever la raison pour laquelle ils sont pleins.
