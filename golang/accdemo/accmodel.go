package main

import (
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
)

var FoodStrs = []string{"Apple", "Banana", "Guava", "Orange", "Blackberry", "Mango", "Kiwi", "Raspberry", "Pineapple", "Avacado", "Onion", "Lettuce", "Cheese", "Almond", "Cake", "Walnut"}
var SuffixStrs = []string{"LLC", "Inc", "Corp", "Corp", "Ltd", "", "", "", "", "", "", "", ""}
var ModWordStrs = []string{"Star", "Lightning", "Flash", "Media", "Data", "Micro", "Net", "Echo", "World", "Red", "Blue", "Green", "Yellow", "Purple", "Tele", "Cloud", "Insta", "Face", "Super"}
var EmailStrs = []string{"mike", "matt", "michelle", "pat", "jim", "marc", "andrew", "alan", "henry", "jenny"}

type AccType struct {
	AccId   string
	AccName string
	IsPaid  bool
	Email   string
}

type AccModel struct {
	Lock *sync.Mutex
	Accs []*AccType
}

func MakeAccModel() *AccModel {
	rtn := &AccModel{Lock: &sync.Mutex{}}
	rtn.RegenAccounts()
	return rtn
}

func (m *AccModel) AccById(id string) *AccType {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	return m.accByIdNoLock(id)
}

func (m *AccModel) CopyAccList() []*AccType {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	rtn := make([]*AccType, len(m.Accs))
	for i := 0; i < len(m.Accs); i++ {
		accCopy := *m.Accs[i]
		rtn[i] = &accCopy
	}
	return rtn
}

func (m *AccModel) accByIdNoLock(id string) *AccType {
	for _, acc := range m.Accs {
		if acc.AccId == id {
			return acc
		}
	}
	return nil
}

func (m *AccModel) CreateAcc(name string, email string) string {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	accId := uuid.New().String()
	m.Accs = append(m.Accs, &AccType{AccId: accId, AccName: name, Email: email})
	return accId
}

func (m *AccModel) RemoveAcc(id string) bool {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	pos := -1
	for idx, acc := range m.Accs {
		if acc.AccId == id {
			pos = idx
			break
		}
	}
	if pos == -1 {
		return false
	}
	m.Accs = append(m.Accs[:pos], m.Accs[pos+1:]...)
	return true
}

func (m *AccModel) Upgrade(id string) bool {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	acc := m.accByIdNoLock(id)
	if acc == nil {
		return false
	}
	acc.IsPaid = true
	return true
}

func (m *AccModel) Downgrade(id string) bool {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	acc := m.accByIdNoLock(id)
	if acc == nil {
		return false
	}
	acc.IsPaid = false
	return true
}

func (m *AccModel) RegenAccounts() {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.Accs = nil
	for i := 0; i < 5; i++ {
		m.Accs = append(m.Accs, MakeRandomAcc())
	}
}

func RandomWord(list []string) string {
	return list[rand.Int31n(int32(len(list)))]
}

func MakeRandomAcc() *AccType {
	return &AccType{
		AccId:   uuid.New().String(),
		AccName: strings.TrimSpace(RandomWord(ModWordStrs) + " " + RandomWord(FoodStrs) + " " + RandomWord(SuffixStrs)),
		IsPaid:  false,
		Email:   RandomWord(EmailStrs) + strconv.Itoa(int(rand.Int31n(70)+10)) + "@nomail.com",
	}
}
