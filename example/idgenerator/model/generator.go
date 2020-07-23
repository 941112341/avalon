package model

import (
	"github.com/941112341/avalon/example/idgenerator/model/repository"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/pkg/errors"
	"sync"
)

const (
	extendBase   = 1000
	requestLimit = 100
)

type IdGeneratorModel struct {
	repository.IdGenerator

	repository repository.IdGeneratorRepository

	Index int64
	IDs   []int64
	lock  sync.Mutex
}

func (i *IdGeneratorModel) nextMaxId(cnt int64) int64 {
	return i.Index + cnt - 1
}

func (i *IdGeneratorModel) remain() int64 {
	return i.Length - i.Index
}

// return nil if cnt out of range
func (i *IdGeneratorModel) subIds(cnt int64) []int64 {
	if i.canAssign(cnt) {
		maxID := i.nextMaxId(cnt)
		ids := i.IDs[i.Index : maxID+1]
		return ids
	}
	return nil
}

func (i *IdGeneratorModel) GetIds() []int64 {
	return i.IDs
}

func (i *IdGeneratorModel) SetIndex(index int64) {
	i.Index = index
}

func (i *IdGeneratorModel) AddIndex(cnt int64) {
	i.Index += cnt
}

func (i *IdGeneratorModel) subIdGenerator(cnt int64, bizID string) (*IdGeneratorModel, error) {
	if err := i.valid(); err != nil {
		return nil, err
	}
	if !i.canAssign(cnt) {
		return nil, errors.New("cannot assign")
	}

	maxID := i.nextMaxId(cnt)
	ids := i.subIds(cnt) // len == cnt

	subModel := &IdGeneratorModel{
		IdGenerator: repository.IdGenerator{
			ID:      0,
			MaxID:   maxID,
			Length:  cnt,
			BizID:   bizID,
			Version: 0,
		},
		repository: i.repository,
		Index:      0,
		IDs:        ids,
		lock:       sync.Mutex{},
	}
	if err := subModel.valid(); err != nil {
		return nil, err
	}
	return subModel, nil
}

func (i *IdGeneratorModel) canAssign(cnt int64) bool {
	if cnt > requestLimit || cnt < 0 {
		return false
	}
	return int64(len(i.IDs))-i.Index >= cnt
}

func (i *IdGeneratorModel) valid() error {
	if i == nil {
		return errors.New("nil ptr")
	}
	if i.Index < 0 {
		return errors.New("index < 0")
	}
	if i.repository == nil {
		return errors.New("repository is nil")
	}
	if len(i.IDs) == 0 {
		return errors.New("ids is nil")
	}
	if *inline.LastInt64(i.IDs) != i.MaxID {
		return errors.New("maxID != last int in ids")
	}
	return nil
}

func (i *IdGeneratorModel) Assign(cnt int64, bizId string) ([]int64, error) {
	i.lock.Lock()
	defer i.lock.Unlock()
	if i.canAssign(cnt) {
		subModel, err := i.subIdGenerator(cnt, bizId)
		if err != nil {
			inline.WithFields("method", "assign", "i", inline.ToJsonString(i), "cnt", cnt).Error("subIdGenerator fail")
			return nil, errors.Wrap(err, "sub generator")
		}

		if err := subModel.save(); err != nil {
			return nil, err
		}
		return subModel.IDs, nil
	}
	modelAfterExtend, err := i.extend()
	if err != nil {
		return nil, errors.Wrap(err, "model after extend")
	}
	return modelAfterExtend.Assign(cnt, bizId)
}

func (i *IdGeneratorModel) save() error {
	err := i.valid()
	if err != nil {
		return errors.Wrap(err, "valid")
	}
	return i.repository.Save(i.IdGenerator)
}

func (i *IdGeneratorModel) extend() (*IdGeneratorModel, error) {

	m, err := i.newModel(extendBase)
	if err != nil {
		return nil, errors.Wrap(err, "new model")
	}
	e := i.merge(m)
	return e, nil
}

func (i *IdGeneratorModel) merge(n *IdGeneratorModel) *IdGeneratorModel {
	if n == nil {
		return i
	}
	return NewGenerator(repository.IdGenerator{
		ID:      n.ID,
		MaxID:   n.MaxID,
		Length:  n.Length + i.remain(),
		BizID:   n.BizID,
		Version: n.Version,
	}, i.repository)
}

func (i *IdGeneratorModel) newModel(cnt int64) (*IdGeneratorModel, error) {

	subModel := NewGenerator(repository.IdGenerator{
		ID:      0,
		MaxID:   i.nextMaxId(cnt),
		Length:  cnt,
		BizID:   i.BizID,
		Version: i.Version,
	}, i.repository)

	err := inline.Retry(func() error {
		rows, err := subModel.repository.UpdateVersion(subModel.IdGenerator)
		if err != nil {
			return err
		}
		if rows == 0 {
			subModel, err = NewGeneratorDB(subModel.repository)
			if err != nil {
				return err
			}
			subModel = NewGenerator(repository.IdGenerator{
				ID:      0,
				MaxID:   subModel.nextMaxId(cnt),
				Length:  cnt,
				BizID:   subModel.BizID,
				Version: subModel.Version,
			}, subModel.repository)
		} else {
			subModel.Version++
		}

		return nil
	}, 3, 0)

	return subModel, err
}

func NewGeneratorDB(repository repository.IdGeneratorRepository) (*IdGeneratorModel, error) {
	generator, err := repository.Get()
	if err != nil {
		return nil, err
	}
	return NewGenerator(*generator, repository), nil
}

func NewGenerator(generator repository.IdGenerator, repository repository.IdGeneratorRepository) *IdGeneratorModel {

	return &IdGeneratorModel{
		IdGenerator: generator,
		repository:  repository,
		Index:       0,
		IDs:         inline.BuildIntList(generator.MaxID, generator.Length),
	}
}