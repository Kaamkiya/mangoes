package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"

	"github.com/charmbracelet/huh"
)

func addNote(tr *smf.Track, note uint8, velocity uint8, duration uint32) {
	tr.Add(0, midi.NoteOn(0, note, velocity))
	tr.Add(duration, midi.NoteOff(0, note))
}

var (
	CmajTriad  = Chord{midi.C(5), midi.E(5), midi.G(5)}
	DbmajTriad = Chord{midi.Db(5), midi.F(5), midi.Ab(5)}
	DmajTriad  = Chord{midi.D(5), midi.Gb(5), midi.A(5)}
	EbmajTriad = Chord{midi.Eb(5), midi.G(5), midi.Bb(5)}
	EmajTriad  = Chord{midi.E(5), midi.Ab(5), midi.B(5)}
	FmajTriad  = Chord{midi.F(5), midi.A(5), midi.C(5)}
	GbmajTriad = Chord{midi.Gb(5), midi.Bb(5), midi.Db(5)}
	GmajTriad  = Chord{midi.G(5), midi.B(5), midi.D(5)}
	AbmajTriad = Chord{midi.Ab(5), midi.C(5), midi.Eb(5)}
	AmajTriad  = Chord{midi.A(5), midi.Db(5), midi.E(5)}
	BbmajTriad = Chord{midi.Bb(5), midi.D(5), midi.F(5)}
	BmajTriad  = Chord{midi.B(5), midi.Eb(5), midi.Gb(5)}

	CminTriad  = Chord{midi.C(5), midi.Eb(5), midi.G(5)}
	CSminTriad = Chord{midi.Db(5), midi.E(5), midi.Ab(5)}
	DminTriad  = Chord{midi.D(5), midi.F(5), midi.A(5)}
	EbminTriad = Chord{midi.Eb(5), midi.Gb(5), midi.Bb(5)}
	EminTriad  = Chord{midi.E(5), midi.G(5), midi.B(5)}
	FminTriad  = Chord{midi.F(5), midi.Ab(5), midi.C(5)}
	FSminTriad = Chord{midi.Gb(5), midi.A(5), midi.Db(5)}
	GminTriad  = Chord{midi.G(5), midi.Bb(5), midi.D(5)}
	GSminTriad = Chord{midi.Ab(5), midi.B(5), midi.Eb(5)}
	AminTriad  = Chord{midi.A(5), midi.C(5), midi.E(5)}
	BbminTriad = Chord{midi.Bb(5), midi.Db(5), midi.F(5)}
	BminTriad  = Chord{midi.B(5), midi.D(5), midi.Gb(5)}
)

type Chord [3]uint8
type ChordProgression [4]Chord

// TODO:
// make it so that users can choose the chord progression. they submit a
// progression, eg I-ii-IV-iii, and a key, eg C, and the progression is
// generated

func getInputs() (tempo float64, bars uint64, progression ChordProgression) {
	var barsStr string
	var tempoStr string

	chordOptions := []huh.Option[Chord]{
		huh.NewOption("C maj", CmajTriad),
		huh.NewOption("Db maj", DbmajTriad),
		huh.NewOption("D maj", DmajTriad),
		huh.NewOption("Eb maj", EbmajTriad),
		huh.NewOption("E maj", EmajTriad),
		huh.NewOption("F maj", FmajTriad),
		huh.NewOption("Gb maj", GbmajTriad),
		huh.NewOption("G maj", GmajTriad),
		huh.NewOption("Ab maj", AbmajTriad),
		huh.NewOption("A maj", AmajTriad),
		huh.NewOption("Bb maj", BbmajTriad),
		huh.NewOption("B maj", BmajTriad),
		huh.NewOption("C min", CminTriad),
		huh.NewOption("C# min", CSminTriad),
		huh.NewOption("D min", DminTriad),
		huh.NewOption("Eb min", EbminTriad),
		huh.NewOption("E min", EminTriad),
		huh.NewOption("F min", FminTriad),
		huh.NewOption("F# min", FSminTriad),
		huh.NewOption("G min", GminTriad),
		huh.NewOption("G# min", GSminTriad),
		huh.NewOption("A min", AminTriad),
		huh.NewOption("Bb min", BbminTriad),
		huh.NewOption("B min", BminTriad),
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter a tempo (BPM)").
				Value(&tempoStr),
			huh.NewInput().
				Title("How many bars of music?").
				Value(&barsStr),
		),
		huh.NewGroup(
			huh.NewSelect[Chord]().
				Title("First triad for the chord progression to use").
				Value(&progression[0]).
				Options(chordOptions...),
		),
		huh.NewGroup(
			huh.NewSelect[Chord]().
				Title("Second triad for the chord progression to use").
				Value(&progression[1]).
				Options(chordOptions...),
		),
		huh.NewGroup(
			huh.NewSelect[Chord]().
				Title("Third triad for the chord progression to use").
				Value(&progression[2]).
				Options(chordOptions...),
		),
		huh.NewGroup(
			huh.NewSelect[Chord]().
				Title("Fourth triad for the chord progression to use").
				Value(&progression[3]).
				Options(chordOptions...),
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

	fmt.Println("cmin-EMAJ-FMAJ-gmin")
	fmt.Println("CMAJ-GMAJ-amin-FMAJ")
	fmt.Println("amin-DMAJ-fmin-CMAJ")

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

		if b%2 == 0 {
			// Play the chord for the bass line
			bass.Add(0, midi.NoteOn(0, chord[0]-24, volume-20))
			bass.Add(0, midi.NoteOn(0, chord[1]-24, volume-20))
			bass.Add(0, midi.NoteOn(0, chord[2]-24, volume-20))
			bass.Add(4*clock.Ticks4th(), midi.NoteOff(0, chord[0]-24))
			bass.Add(0, midi.NoteOff(0, chord[1]-24))
			bass.Add(0, midi.NoteOff(0, chord[2]-24))
		} else {
			// Play the chord twice, as half notes
			bass.Add(0, midi.NoteOn(0, chord[0]-24, volume-20))
			bass.Add(0, midi.NoteOn(0, chord[1]-24, volume-20))
			bass.Add(0, midi.NoteOn(0, chord[2]-24, volume-20))
			bass.Add(2*clock.Ticks4th(), midi.NoteOff(0, chord[0]-24))
			bass.Add(0, midi.NoteOff(0, chord[1]-24))
			bass.Add(0, midi.NoteOff(0, chord[2]-24))

			bass.Add(0, midi.NoteOn(0, chord[0]-12, volume-20))
			bass.Add(0, midi.NoteOn(0, chord[1]-12, volume-20))
			bass.Add(0, midi.NoteOn(0, chord[2]-12, volume-20))
			bass.Add(2*clock.Ticks4th(), midi.NoteOff(0, chord[0]-12))
			bass.Add(0, midi.NoteOff(0, chord[1]-12))
			bass.Add(0, midi.NoteOff(0, chord[2]-12))
		}

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
