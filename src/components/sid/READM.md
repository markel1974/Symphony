# SID MOS 6581 Register Map

The SID MOS 6581 has 29 registers mapped in the Commodore 64's I/O address space, at addresses from `$D400` to `$D41C` (54272 to 54300 in decimal). Here's the list of registers, sorted by address:

## Registers by Address

*   **`$D400` (54272) - `FREQ_LO_1` (Frequency Low 1):** Least significant byte of the frequency for voice 1.
*   **`$D401` (54273) - `FREQ_HI_1` (Frequency High 1):** Most significant byte of the frequency for voice 1.
*   **`$D402` (54274) - `PW_LO_1` (Pulse Width Low 1):** Least significant byte of the pulse width for voice 1.
*   **`$D403` (54275) - `PW_HI_1` (Pulse Width High 1):** Most significant byte of the pulse width for voice 1.
*   **`$D404` (54276) - `CR_1` (Control Register 1):** Control register for voice 1 (waveform, ring modulation, sync, gate).
*   **`$D405` (54277) - `AD_1` (Attack/Decay 1):** Attack and Decay for voice 1.
*   **`$D406` (54278) - `SR_1` (Sustain/Release 1):** Sustain and Release for voice 1.
*   **`$D407` (54279) - `FREQ_LO_2` (Frequency Low 2):** Least significant byte of the frequency for voice 2.
*   **`$D408` (54280) - `FREQ_HI_2` (Frequency High 2):** Most significant byte of the frequency for voice 2.
*   **`$D409` (54281) - `PW_LO_2` (Pulse Width Low 2):** Least significant byte of the pulse width for voice 2.
*   **`$D40A` (54282) - `PW_HI_2` (Pulse Width High 2):** Most significant byte of the pulse width for voice 2.
*   **`$D40B` (54283) - `CR_2` (Control Register 2):** Control register for voice 2 (waveform, ring modulation, sync, gate).
*   **`$D40C` (54284) - `AD_2` (Attack/Decay 2):** Attack and Decay for voice 2.
*   **`$D40D` (54285) - `SR_2` (Sustain/Release 2):** Sustain and Release for voice 2.
*   **`$D40E` (54286) - `FREQ_LO_3` (Frequency Low 3):** Least significant byte of the frequency for voice 3.
*   **`$D40F` (54287) - `FREQ_HI_3` (Frequency High 3):** Most significant byte of the frequency for voice 3.
*   **`$D410` (54288) - `PW_LO_3` (Pulse Width Low 3):** Least significant byte of the pulse width for voice 3.
*   **`$D411` (54289) - `PW_HI_3` (Pulse Width High 3):** Most significant byte of the pulse width for voice 3.
*   **`$D412` (54290) - `CR_3` (Control Register 3):** Control register for voice 3 (waveform, ring modulation, sync, gate).
*   **`$D413` (54291) - `AD_3` (Attack/Decay 3):** Attack and Decay for voice 3.
*   **`$D414` (54292) - `SR_3` (Sustain/Release 3):** Sustain and Release for voice 3.
*   **`$D415` (54293) - `FC_LO` (Filter Cutoff Frequency Low):** Least significant byte of the filter cutoff frequency.
*   **`$D416` (54294) - `FC_HI` (Filter Cutoff Frequency High):** Most significant byte of the filter cutoff frequency.
*   **`$D417` (54295) - `RES_FT_VOL` (Filter Resonance/Filter Type/Volume):** Filter resonance, filter type, and master volume.
*   **`$D419` (54297) - `POT_X` (Potentiometer X):** Potentiometer X value.
*   **`$D41A` (54298) - `POT_Y` (Potentiometer Y):** Potentiometer Y value.
*   **`$D41B` (54299) - `OSC3_RANDOM` (Oscillator 3/Random):** Reading of oscillator 3 or a random number.
*   **`$D41C` (54300) - `EF_3` (Envelope Flag/Voice 3):** Envelope flag, Voice 3.

## Register Grouping

For a better understanding, these registers can be grouped by function:

*   **Voice 1 Registers:** `$D400` - `$D406`
*   **Voice 2 Registers:** `$D407` - `$D40D`
*   **Voice 3 Registers:** `$D40E` - `$D414`
*   **Filter/Volume Registers:** `$D415` - `$D417`
* **Potentiometer/Miscellaneous Registers:** `$D419` - `$D41C`
* **Input Registers:**`$D419` - `$D41A`

## Notes

* Each voice has its own frequency (`FREQ_LO_n`, `FREQ_HI_n`), pulse width (`PW_LO_n`, `PW_HI_n`), and control (`CR_n`) registers, as well as `AD` and `SR` registers.
* The registers `POT_X` and `POT_Y` are used for managing the potentiometers.
* The `OSC3_RANDOM` register can be used to read the value of the oscillator 3 or to get a random number.
* The `EF_3` is used for the envelope flag of the voice 3.
* `CR` stands for "Control Register."
* `AD` and `SR` stand for "Attack/Decay" and "Sustain/Release," respectively.
* `PW` and `FREQ` stand for "Pulse Width" and "Frequency," respectively.
* `FC` stands for "Filter Cutoff."
* `RES_FT_VOL` is the register of the filter, the type of filter, and the volume.
* `EF` stands for envelope flag.