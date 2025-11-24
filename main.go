package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"

	"github.com/charmbracelet/huh"
)

// TODO:
// add ABAB form for bass line (bassForm variable decides whether form=chords || form=ABAB)

func addNote(tr *smf.Track, note uint8, velocity uint8, duration uint32) {
	tr.Add(0, midi.NoteOn(0, note, velocity))
	tr.Add(duration, midi.NoteOff(0, note))
}

type Chord [3]uint8
type ChordProgression [4]Chord

var (
	// C F# Bb C is heartbreaking. C Eb7 G#dim7 Amaj9#11 is extremely sorrowful too.

	// Happy, use with 120-150 BPM
	cgaf = ChordProgression{
		Chord{midi.C(5), midi.E(5), midi.G(5)}, // Cmaj
		Chord{midi.G(5), midi.B(5), midi.D(5)}, // Gmaj
		Chord{midi.A(5), midi.C(5), midi.E(5)}, // Amin
		Chord{midi.F(5), midi.A(5), midi.C(5)}, // Fmaj
	}

	// Happy, use with 120-150 BPM
	cfgf = ChordProgression{
		Chord{midi.C(5), midi.E(5), midi.G(5)}, // Cmaj
		Chord{midi.F(5), midi.A(5), midi.C(5)}, // Fmaj
		Chord{midi.G(5), midi.B(5), midi.D(5)}, // Gmaj
		Chord{midi.F(5), midi.A(5), midi.C(5)}, // Fmaj
	}

	// Creepy, use with 80-100 BPM
	cece = ChordProgression{
		Chord{midi.C(4), midi.Eb(4), midi.G(4)},   // Cmaj
		Chord{midi.Eb(4), midi.Gb(4), midi.Bb(4)}, // Ebmaj
		Chord{midi.C(4), midi.Eb(4), midi.G(4)},   // Cmaj
		Chord{midi.Eb(4), midi.Gb(4), midi.Bb(4)}, // Ebmaj
	}

	// Sad, use with 60-77 BPM
	amdfmc = ChordProgression{
		Chord{midi.A(4), midi.C(4), midi.E(4)},  // Amin
		Chord{midi.D(4), midi.Gb(4), midi.A(4)}, // Dmaj
		Chord{midi.F(4), midi.Ab(4), midi.C(4)}, // Fmin
		Chord{midi.C(4), midi.E(4), midi.G(4)},  // Cmaj
	}
)

// TODO:
// make it so that users can choose the chord progression. they submit a
// progression, eg I-ii-IV-iii, and a key, eg C, and the progression is
// generated

func getInputs() (tempo float64, bars uint64, progression ChordProgression) {
	var barsStr string
	var tempoStr string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter a tempo (BPM)").
				Value(&tempoStr),
			huh.NewInput().
				Title("How many bars of music?").
				Value(&barsStr),
			huh.NewSelect[ChordProgression]().
				Title("Chord progression to use").
				Value(&progression).
				Options(
					huh.NewOption("C, G, A, F", cgaf),
					huh.NewOption("C, Eb, C, Eb", cece),
					huh.NewOption("C, F, G, F", cfgf),
					huh.NewOption("Am, D, Fm, C", amdfmc),
				),
		),
	)

	err := form.Run()
	if err != nil {
		panic(err)
	}

	tempo, err = strconv.ParseFloat(tempoStr, 64)
	if err != nil {
		panic(err)
	}

	bars, err = strconv.ParseUint(barsStr, 10, 32)
	if err != nil {
		panic(err)
	}

	return
}

func main() {
	volume := uint8(100)

	tempo, bars, progressionToUse := getInputs()

	melody := smf.Track{}
	bass := smf.Track{}
	s := smf.New()
	clock := smf.MetricTicks(96)

	melody.Add(0, smf.MetaMeter(4, 4))
	melody.Add(0, smf.MetaTempo(tempo))
	bass.Add(0, smf.MetaMeter(4, 4))
	bass.Add(0, smf.MetaTempo(tempo))

	progressionIdx := 0
	var lastRestPosition uint64
	for b := range bars {
		chord := progressionToUse[progressionIdx]

		for i := range 4 {
			note := chord[rand.IntN(len(chord))]
			noteVolume := volume

			shouldPlay8thNotePair := rand.IntN(4) == 0
			canPlayRest := b-lastRestPosition > 3
			shouldPlayRest := rand.IntN(15) == 0
			shouldBeAccented := rand.IntN(10) == 0

			if i%2 == 0 && shouldBeAccented {
				noteVolume = 127
			}

			if shouldPlay8thNotePair { // 1/4 chance of a pair of eighth notes
				addNote(&melody, note, noteVolume, clock.Ticks8th())
				secondNote := chord[rand.IntN(len(chord))]
				for secondNote-note > 6 {
					secondNote = chord[rand.IntN(len(chord))]
				}
				addNote(&melody, secondNote, noteVolume, clock.Ticks8th())
			} else if canPlayRest && shouldPlayRest { // 1/15 chance of a quarter rest
				lastRestPosition = b
				melody.Add(clock.Ticks4th(), midi.NoteOff(0, 0))
			} else {
				addNote(&melody, note, noteVolume, clock.Ticks4th())
			}
		} // for i := range 4

		// Play the chord for the bass line
		bass.Add(0, midi.NoteOn(0, chord[0]-24, volume-20))
		bass.Add(0, midi.NoteOn(0, chord[1]-24, volume-20))
		bass.Add(0, midi.NoteOn(0, chord[2]-24, volume-20))
		bass.Add(4*clock.Ticks4th(), midi.NoteOff(0, chord[0]-24))
		bass.Add(0, midi.NoteOff(0, chord[1]-24))
		bass.Add(0, midi.NoteOff(0, chord[2]-24))

		fmt.Println(chord)

		progressionIdx++
		if progressionIdx > 3 {
			progressionIdx = 0
		}
	} // for b := range bars

	melody.Close(0)
	bass.Close(0)

	s.TimeFormat = smf.MetricTicks(96)
	s.Add(melody)
	s.Add(bass)
	s.WriteFile("mangoes.mid")
}
