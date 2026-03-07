; Nintendo Logo Display Example
; This shows how to decompress and display the Nintendo logo

; The compressed logo data (from ROM header 0x0104-0x0133)
NintendoLogoData:
    DB $CE, $ED, $66, $66, $CC, $0D, $00, $0B
    DB $03, $73, $00, $83, $00, $0C, $00, $0D
    DB $00, $08, $11, $1F, $88, $89, $00, $0E
    DB $DC, $CC, $6E, $E6, $DD, $DD, $D9, $99
    DB $BB, $BB, $67, $63, $6E, $0E, $EC, $CC
    DB $DD, $DC, $99, $9F, $BB, $B9, $33, $3E

; Decompression routine
; This converts the 48 compressed bytes into 384 bytes of tile data
; (48 bytes × 8 bits per byte = 384 bits → 48 tiles of 8 pixels each)

DecompressNintendoLogo:
    ; Source: NintendoLogoData (48 bytes)
    ; Destination: VRAM tile data starting at $8010 (tiles 1-12)
    
    ld hl, NintendoLogoData     ; Source pointer
    ld de, $8010                ; Destination in VRAM (tile 1)
    ld b, 48                    ; 48 bytes to decompress
    
.decompress_loop:
    ld a, [hl+]                 ; Read compressed byte
    ld c, a                     ; Save it in C
    
    ; Each compressed byte becomes 16 bytes of tile data
    ; (2 bytes per row × 8 rows per tile, but spread across 2 tiles)
    
    ; Here's how decompression works:
    ; Each bit in the compressed byte represents 2 bits in output
    ; Bit=1 becomes "11", Bit=0 becomes "00"
    
    ; Process 8 bits (will create 2 bytes of output per bit)
    ld a, c
    call ExpandByte
    
    dec b
    jr nz, .decompress_loop
    ret

; Helper: Expand one byte into 2 bytes of tile data
; This is how the Game Boy boot ROM actually does it:
; Each bit becomes 2 consecutive bits in the output
;
; Example: 0xCE = %11001110
; Becomes: %11110000 %11111100
;
; Input: A = compressed byte  
; Output: 2 bytes written to [DE]
ExpandByte:
    ld c, a                     ; Save original byte
    
    ; First output byte: bits 7,7,6,6,5,5,4,4
    ld b, 4                     ; Process 4 bit pairs
    xor a                       ; Clear A
    
.first_byte:
    rl c                        ; Get next bit into carry
    rla                         ; Put it in A
    rl c                        ; Same bit again into carry  
    rla                         ; Put it in A again
    dec b
    jr nz, .first_byte
    
    ld [de], a                  ; Write first byte
    inc de
    
    ; Second output byte: bits 3,3,2,2,1,1,0,0
    ld b, 4                     ; Process 4 more bit pairs
    xor a                       ; Clear A
    
.second_byte:
    rl c                        ; Get next bit into carry
    rla                         ; Put it in A
    rl c                        ; Same bit again into carry
    rla                         ; Put it in A again  
    dec b
    jr nz, .second_byte
    
    ld [de], a                  ; Write second byte
    inc de
    ret

; To actually DISPLAY the logo:
; You need to set up a tilemap that references these tiles

DisplayNintendoLogo:
    ; First, decompress the logo into VRAM
    call DecompressNintendoLogo
    
    ; Now create a tilemap
    ; The logo is 12 tiles wide × 2 tiles tall = 24 tiles total
    ; We'll display it centered on screen
    
    ld hl, $9800 + 32*6 + 4    ; Tilemap address (row 6, col 4)
    ld b, 12                    ; 12 tiles across
    ld a, 1                     ; Starting tile number
    
.draw_top_row:
    ld [hl+], a                 ; Write tile number to tilemap
    inc a
    dec b
    jr nz, .draw_top_row
    
    ; Draw bottom row
    ld hl, $9800 + 32*7 + 4    ; Next row down
    ld b, 12
    
.draw_bottom_row:
    ld [hl+], a
    inc a
    dec b
    jr nz, .draw_bottom_row
    
    ret
