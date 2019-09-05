package main

import "github.com/asdine/storm"

type Muted struct {
	UserId int `storm:"id"`
}

func (m Muted) Mute() error {
	if err := storage.Save(&m); err != nil {
		log.Debugf("Failed to save data to storage, error: %s", err)
		return err
	}
	return nil
}

func (m Muted) Unmute() error {
	if err := storage.DeleteStruct(&m); (err != nil) && (err != storm.ErrNotFound) {
		log.Debugf("Failed to remove data from storage, error: %s", err)
		return err
	}
	return nil
}

func (m Muted) IsMuted() (bool, error) {
	err := storage.One("UserId", m.UserId, &m)
	if err == storm.ErrNotFound {
		return false, nil
	} else if err != nil {
		log.Debugf("Can not get data from storage, error: %s", err)
		return false, err
	}
	return true, nil
}

type LastChanMessage struct {
	ChatID    int64 `storm:"id"`
	MessageID int
	Sender    int
}

func (m LastChanMessage) Save() error {
	if err := storage.Save(&m); err != nil {
		log.Debugf("Can not save last messageID, error: %s", err)
		return err
	}
	return nil
}

type Integration struct {
	ChatID int64 `storm:"id"`
}

func (i Integration) Add() error {
	if err := storage.Save(&i); err != nil {
		log.Debugf("Can not add integration, error: %s", err)
		return err
	}
	return nil
}

func (i Integration) Remove() error {
	if err := storage.DeleteStruct(&i); err != nil {
		log.Debugf("Can not add integration, error: %s", err)
		return err
	}
	return nil
}

func (i Integration) GetAll() ([]Integration, error) {
	var itg []Integration
	if err := storage.All(&itg); err != nil {
		log.Debugf("Can not get integrations, error: %s", err)
		return nil, err
	}
	return itg, nil
}
