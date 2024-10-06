package service

import "github.com/sirupsen/logrus"

type GroupsRepo interface {
}

type SongsRepo interface {
}

type ApplicationServices struct {
	log    *logrus.Logger
	groups GroupsRepo
	songs  SongsRepo
}

func New(log *logrus.Logger, groups GroupsRepo, songs SongsRepo) *ApplicationServices {
	return &ApplicationServices{
		log:    log,
		groups: groups,
		songs:  songs,
	}
}
