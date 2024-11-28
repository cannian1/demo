package main

import (
	"errors"
	"fmt"
)

type TVState interface {
	On() error
	Off() error
	Mute() error
}

type OnState struct{}

func (s *OnState) On() error {
	return errors.New("TV is already on")
}

func (s *OnState) Off() error {
	fmt.Println("Turning TV off")
	return nil
}

func (s *OnState) Mute() error {
	fmt.Println("Muting TV")
	return nil
}

type OffState struct{}

func (s *OffState) On() error {
	fmt.Println("Turning TV on")
	return nil
}

func (s *OffState) Off() error {
	return errors.New("TV is already off")
}

func (s *OffState) Mute() error {
	return errors.New("cannot mute TV when it's off")
}

type MuteState struct{}

func (s *MuteState) On() error {
	fmt.Println("Turning TV on")
	return nil
}

func (s *MuteState) Off() error {
	fmt.Println("Turning TV off")
	return nil
}

func (s *MuteState) Mute() error {
	return errors.New("TV is already muted")
}

type TV struct {
	state TVState
}

func (t *TV) setState(state TVState) {
	t.state = state
}

func (t *TV) On() error {
	if err := t.state.On(); err != nil {
		return err
	}
	t.setState(new(OnState))

	return nil
}

func (t *TV) Off() error {
	if err := t.state.Off(); err != nil {
		return err
	}
	t.setState(new(OffState))

	return nil
}

func (t *TV) Mute() error {
	if err := t.state.Mute(); err != nil {
		return err
	}
	t.setState(new(MuteState))

	return nil
}

func main() {
	tv := &TV{&OffState{}}

	err := tv.On()
	if err != nil {
		fmt.Println(err)
	}

	err = tv.Mute()
	if err != nil {
		fmt.Println(err)
	}

	err = tv.Off()
	if err != nil {
		fmt.Println(err)
	}

	err = tv.Mute()
	if err != nil {
		fmt.Println(err)
	}
}
