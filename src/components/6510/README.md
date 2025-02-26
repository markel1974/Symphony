# Package mos6510

Questo package implementa l'emulazione del microprocessore MOS 6510 (una variante del 6502) utilizzato nel Commodore 64.

## Panoramica

Il package `mos6510` fornisce una rappresentazione software della CPU 6510, inclusi:

*   Registri (A, X, Y, PC, SP, SR).
*   Flag del processore (N, V, B, D, I, Z, C).
*   Ciclo di esecuzione delle istruzioni (fetch, decode, execute).
*   Gestione degli interrupt (NMI, IRQ, Reset).
*   Implementazione delle istruzioni del 6510 (tramite micro-operazioni).
*   Gestione dello stack.
* Tabelle di lookup

## Struttura del Package

Il package è organizzato nei seguenti file:

*   `cpu.go`: Definisce la struct `CPU` e i metodi principali per l'emulazione (ciclo di esecuzione, gestione registri, ecc.).
*   `instructions.go`: Contiene le *dichiarazioni* delle funzioni che implementano le singole istruzioni del 6510 (suddivise in micro-operazioni).
*   `inst_*.go`: Contengono l'*implementazione* delle micro-operazioni delle istruzioni, raggruppate per categoria (load/store, aritmetiche, logiche, ecc.).
*   `opcodes.go`: Definisce le tabelle di dispatch (`_modeTable` e `_opTable`) che mappano gli opcode alle funzioni di gestione delle modalità di indirizzamento e alle funzioni di esecuzione delle istruzioni.
*   `stack.go`: Implementa le operazioni sullo stack del 6510.
*   `utils.go`: Contiene funzioni di utilità.
*   `interrupts_test.go`: Contiene i test per la gestione degli interrupt.
* `opcodes_test.go`: Contiene test per le operations.

## Istruzioni Implementate

[**TODO:** Elencare *tutte* le istruzioni implementate, con una breve descrizione di ciascuna, la modalità di indirizzamento, i flag modificati, e i cicli di clock.  Questo può essere fatto in forma di tabella, o usando una lista.]

**Esempio:**

| Istruzione | Modalità di Indirizzamento | Descrizione                                   | Flag Affetti | Cicli |
| :---------- | :------------------------- | :-------------------------------------------- | :----------- | ----- |
| LDA         | Immediate                  | Carica un valore immediato nell'accumulatore. | N, Z         | 2     |
| LDA         | Zero Page                  | Carica un valore da un indirizzo in Zero Page.  | N, Z         | 3     |
| ...         | ...                        | ...                                           | ...          | ...   |

## Modalità di Indirizzamento

[**TODO:** Descrivere le modalità di indirizzamento del 6502/6510, con esempi.]

## Interrupt

[**TODO:** Spiegare come vengono gestiti gli interrupt (NMI, IRQ, Reset).]

## Dipendenze

*   `github.com/markel1974/c64emu/src/memory` (per l'accesso alla memoria)
*   `github.com/markel1974/c64emu/src/components/quartz` (per la gestione del clock)
* Altre interfacce

## Note

*   Questo emulatore implementa *tutte* le istruzioni non documentate del 6502/6510.
*   Questo emulatore *mira* all'accuratezza ciclo per ciclo.

## TODO

*   Aggiungere test unitari per *tutte* le istruzioni, in *tutte* le modalità di indirizzamento.
*   Migliorare la gestione degli errori.
*   Aggiungere commenti dettagliati alle micro-operazioni.
*   Completare l'implementazione delle istruzioni mancanti (se ce ne sono).

## Contribuire

[**TODO:** Se accetti contributi, spiega come farlo.]

## Licenza

Questo progetto è rilasciato sotto licenza [Apache 2.0](https://opensource.org/licenses/Apache-2.0).