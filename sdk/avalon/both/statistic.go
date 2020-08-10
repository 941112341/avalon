package both

import (
	"errors"
	_map "github.com/941112341/avalon/sdk/collect/map"
	"github.com/941112341/avalon/sdk/inline"
	"time"
)

type Statistic interface {
	CountErr(key string) int
	CountSuc(key string) int
	StatisticResult(key string) StatisticResult
	ErrRate(key string) float64
	AddErr(key string)
	AddSuc(key string)
}

type RecordStatistic interface {
	Statistic

	StatisticResultBetween(key string, start, end time.Time) StatisticResult
	CountErrBetween(key string, start, end time.Time) int
	CountSucBetween(key string, start, end time.Time) int
	CountRateBetween(key string, start, end time.Time) float64
	LastStatisticLog(key string, last int) []StatisticLog
	AddErrLog(key string, input interface{}, err error)
	AddSucLog(key string, input, output interface{})
}

type MemoryStatistic struct {
	mapList _map.LimitMap
}

func (m *MemoryStatistic) AddErr(key string) {
	m.AddErrLog(key, nil, errors.New("unknown err"))
}

func (m *MemoryStatistic) AddSuc(key string) {
	m.AddSucLog(key, nil, nil)
}

func (m *MemoryStatistic) AddErrLog(key string, input interface{}, err error) {
	baseLog := &BaseStatistic{
		logTime: time.Now(),
		err:     err,
		input:   input,
		output:  nil,
	}
	m.mapList.Append(key, baseLog)
}

func (m *MemoryStatistic) AddSucLog(key string, input, output interface{}) {

	baseLog := &BaseStatistic{
		logTime: time.Now(),
		err:     nil,
		input:   input,
		output:  output,
	}
	m.mapList.Append(key, baseLog)

}

func (m *MemoryStatistic) StatisticResultBetween(key string, start, end time.Time) StatisticResult {
	result := &BaseStatisticResult{}
	list := m.getListOrInit(key)
	var i, total int
	for _, v := range list {
		if !inline.Between(v.Time(), start, end) {
			continue
		}
		if v.IsError() {
			i++
		}
		total++
	}

	result.sucCnt = total
	result.errCnt = i
	result.sucCnt = total - i
	if total != 0 {
		result.rate = float64(i) / float64(total)
	}
	return result
}

func (m *MemoryStatistic) StatisticResult(key string) StatisticResult {
	return m.StatisticResultBetween(key, time.Time{}, time.Now())
}

func NewMemoryStatistic() *MemoryStatistic {
	return &MemoryStatistic{mapList: _map.NewLimitMap(10)}
}

func (m *MemoryStatistic) getListOrInit(key string) []StatisticLog {
	i := m.mapList.GetOrSet(key, []StatisticLog{})
	log := make([]StatisticLog, 0)
	for _, statisticLog := range i.([]StatisticLog) {
		log = append(log, statisticLog)
	}
	return log
}

func (m *MemoryStatistic) CountErr(key string) int {
	return m.StatisticResult(key).ErrCnt()
}

func (m *MemoryStatistic) CountSuc(key string) int {
	return m.StatisticResult(key).SucCnt()
}

func (m *MemoryStatistic) ErrRate(key string) float64 {
	return m.StatisticResult(key).Rate()
}

func (m *MemoryStatistic) CountErrBetween(key string, start, end time.Time) int {
	return m.StatisticResultBetween(key, start, end).ErrCnt()
}

func (m *MemoryStatistic) CountSucBetween(key string, start, end time.Time) int {
	return m.StatisticResultBetween(key, start, end).SucCnt()
}

func (m *MemoryStatistic) CountRateBetween(key string, start, end time.Time) float64 {
	return m.StatisticResultBetween(key, start, end).Rate()
}

func (m *MemoryStatistic) LastStatisticLog(key string, last int) []StatisticLog {
	list := m.getListOrInit(key)
	idx := len(list) - last
	if idx < 0 {
		idx = 0
	}
	if idx >= len(list) {
		return []StatisticLog{}
	}
	return list[idx:]
}

type StatisticResult interface {
	Total() int
	ErrCnt() int
	SucCnt() int
	Rate() float64

	// other
}

type BaseStatisticResult struct {
	total  int
	errCnt int
	sucCnt int
	rate   float64
}

func (b BaseStatisticResult) Total() int {
	return b.total
}

func (b BaseStatisticResult) ErrCnt() int {
	return b.errCnt
}

func (b BaseStatisticResult) SucCnt() int {
	return b.sucCnt
}

func (b BaseStatisticResult) Rate() float64 {
	return b.rate
}

type StatisticLog interface {
	Time() time.Time
	Error() error
	IsError() bool
	Input() interface{}
	Output() interface{}
}

type BaseStatistic struct {
	logTime       time.Time
	err           error
	input, output interface{}
}

func (b *BaseStatistic) Time() time.Time {
	return b.logTime
}

func (b *BaseStatistic) Error() error {
	return b.err
}

func (b *BaseStatistic) IsError() bool {
	return b.err != nil
}

func (b *BaseStatistic) Input() interface{} {
	return b.input
}

func (b *BaseStatistic) Output() interface{} {
	return b.output
}

type Threshold interface {
	Duration() time.Duration
	Threshold(result StatisticResult) bool
}

type BaseThreshold struct {
	LastDuration time.Duration

	Rate  float64
	Total int
}

func (b *BaseThreshold) Duration() time.Duration {
	return b.LastDuration
}

func (b *BaseThreshold) Threshold(result StatisticResult) bool {
	return result.Rate() > b.Rate && result.ErrCnt() > result.ErrCnt()
}
