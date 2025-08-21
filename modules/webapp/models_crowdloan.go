package webapp

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/jinzhu/gorm"
	"ucoi4tzlito52agmquc6oopn3zpmr6djz5vvfabtgrhyc6hufpzjtnad.onion/Tochka/tochka-free-market/modules/apis"
)

/*
	Models
*/

type Crowdloan struct {
	CloseDate           *time.Time `json:"close_date"`
	ClosingAmount       *float64   `json:"closing_amount"`
	GoalDate            time.Time  `json:"goal_date"`
	GoalValue           float64    `json:"goal_value"`
	PolkadotCrowdloanId uint64     `json:"polkadot_crowdloan_id"`
	Type                string     `json:"type" sql:"index"`
	UserUuid            string     `json:"-" sql:"index"`
	Uuid                string     `json:"uuid" gorm:"primary_key"`
	Value               float64    `json:"value"`
	WeeklyInterestRate  uint16     `json:"weekly_interest_rate"`

	Status []CrowdloanStatus `json:"transaction_status"`
	User   User              `json:"user" gorm:"foreignkey:user_uuid"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CrowdloanStatus struct {
	gorm.Model
	Time          time.Time `json:"time" sql:"index"`
	Status        string    `json:"status" sql:"index"`
	Comment       string    `json:"comment"`
	CrowdloanUuid string    `json:"-" sql:"index"`
}

type Crowdloans []Crowdloan

/*
	View Models
*/

type ViewCrowdloan struct {
	*Crowdloan
	IsFunding       bool `json:"is_funding,omitempty"`
	IsFunded        bool `json:"is_funded,omitempty"`
	IsRefunded      bool `json:"is_refunded,omitempty"`
	IsInterestYield bool `json:"is_interest_yield,omitempty"`

	CurrentStatus     string `json:"current_status,omitempty"`
	CurrentStatusTime string `json:"current_status_time_string,omitempty"`
	CreatedAt         string `json:"created_at_string,omitempty"`
	GoalDate          string `json:"goal_date_string,omitempty"`
}

type ViewCrowdloanLend struct {
	*ViewCrowdloan
	LendValue float64
}

type ViewCrowdloans []ViewCrowdloan

func (loan Crowdloan) ViewCrowdloan(lang string) ViewCrowdloan {
	viewCrowdloan := ViewCrowdloan{
		Crowdloan:         &loan,
		IsFunding:         loan.IsFunding(),
		IsFunded:          loan.IsFunded(),
		IsRefunded:        loan.IsRefunded(),
		IsInterestYield:   loan.IsInterestYield(),
		CurrentStatus:     loan.CurrentStatus(),
		CurrentStatusTime: loan.CurrentStatusTime().Format("02.01.2006 15:04"),
		GoalDate:          loan.GoalDate.Format("02.01.2006 15:04"),
	}

	if loan.CreatedAt != nil {
		viewCrowdloan.CreatedAt = loan.CreatedAt.Format("02.01.2006 15:04")
	}

	return viewCrowdloan
}

func (loans Crowdloans) ViewCrowdloans(lang string) ViewCrowdloans {
	viewLoans := ViewCrowdloans{}
	for _, loan := range loans {
		viewLoans = append(viewLoans, loan.ViewCrowdloan(lang))
	}

	return viewLoans
}

/*
	Status Helpers
*/

func (c Crowdloan) CurrentStatus() string {
	status := c.StatusSortedByTime()
	if len(status) == 0 {
		return "NEW"
	}
	return status[len(status)-1].Status
}

func (c Crowdloan) CurrentStatusTime() time.Time {
	status := c.StatusSortedByTime()
	if len(status) == 0 {
		return time.Now()
	}
	return status[len(status)-1].Time
}

func (c Crowdloan) IsFunding() bool       { return c.CurrentStatus() == "FUNDING" }
func (c Crowdloan) IsFunded() bool        { return c.CurrentStatus() == "FUNDED" }
func (c Crowdloan) IsRefunded() bool      { return c.CurrentStatus() == "REFUNDED" }
func (c Crowdloan) IsInterestYield() bool { return c.CurrentStatus() == "INTEREST_YIELD" }

/*
	Database methods
*/

func (c Crowdloan) Validate() error {
	if c.Uuid == "" {
		return errors.New("wrong uuid")
	}
	if c.Type == "" {
		return errors.New("wrong type")
	}
	if c.UserUuid == "" {
		return errors.New("wrong user uuid")
	}

	if c.GoalValue == 0 {
		return errors.New("wrong goal value")
	}

	if c.GoalDate.Unix() < time.Now().Unix() {
		return errors.New("wrong goal date")
	}

	return nil
}

func (c Crowdloan) Save() error {
	err := c.Validate()
	if err != nil {
		return err
	}
	return c.SaveToDatabase()
}

func (c Crowdloan) Remove() error {
	return database.Delete(c).Error
}

func (c Crowdloan) SaveToDatabase() error {
	if existing, _ := FindCrowdloanByUuid(c.Uuid); existing == nil {
		return database.Create(&c).Error
	}
	return database.Save(&c).Error
}

//

func (c CrowdloanStatus) Validate() error {
	if c.CrowdloanUuid == "" {
		return errors.New("wrong crowdloan uuid")
	}
	if c.Status == "" {
		return errors.New("wrong status")
	}
	return nil
}

func (c CrowdloanStatus) Save() error {
	err := c.Validate()
	if err != nil {
		return err
	}
	return c.SaveToDatabase()
}

func (c CrowdloanStatus) SaveToDatabase() error {
	return database.Create(&c).Error
}

func (t Crowdloan) StatusSortedByTime() []CrowdloanStatus {
	status := t.Status
	sort.Slice(status, func(i, j int) bool {
		return status[i].CreatedAt.Before(status[j].CreatedAt)
	})
	return status
}

/*
	Database queries
*/

func FindCrowdloanByUuid(uuid string) (*Crowdloan, error) {
	var crowdloan Crowdloan
	err := database.
		Preload("Status").
		Preload("User").
		First(&crowdloan, "uuid = ?", uuid).
		Error
	return &crowdloan, err
}

func CreateCrowdloan(
	uuid string,
	crowdloanType string,
	goalValue float64,
	goalDate time.Time,
	weeklyInterestRate uint16,
	userUuid string,
	polkadotCrowdloanId uint64,
) (Crowdloan, error) {
	now := time.Now()
	crowdloan := Crowdloan{
		Uuid:                uuid,
		PolkadotCrowdloanId: polkadotCrowdloanId,
		Type:                crowdloanType,
		GoalValue:           goalValue,
		GoalDate:            goalDate,
		WeeklyInterestRate:  weeklyInterestRate,
		UserUuid:            userUuid,
		Status:              []CrowdloanStatus{},
		CreatedAt:           &now,
	}
	return crowdloan, crowdloan.Save()
}

func ConvertPolkadotCrowdloanFromBlockchainToModel(blockchainCrowdloan apis.PolkdadotCrowdloan, user User) Crowdloan {
	crowdloandUuid := fmt.Sprintf("polkadot-%d", blockchainCrowdloan.Id)
	now := time.Now()
	loan := Crowdloan{
		CreatedAt:           &now,
		GoalDate:            time.Unix(int64(blockchainCrowdloan.GoalDate), 0),
		GoalValue:           blockchainCrowdloan.GoalValue,
		PolkadotCrowdloanId: blockchainCrowdloan.Id,
		Type:                "polkadot",
		User:                user,
		UserUuid:            user.Uuid,
		Uuid:                fmt.Sprintf("polkadot-%d", blockchainCrowdloan.Id),
		Value:               blockchainCrowdloan.Value,
		WeeklyInterestRate:  uint16(blockchainCrowdloan.WeeklyIterest),
		Status: []CrowdloanStatus{
			{
				Model:         gorm.Model{},
				Time:          time.Unix(int64(blockchainCrowdloan.EscrowStatus.Time/1000), 0),
				Status:        blockchainCrowdloan.EscrowStatus.Status,
				Comment:       "",
				CrowdloanUuid: crowdloandUuid,
			},
		},
	}

	if blockchainCrowdloan.CloseDate != nil {
		closeDate := time.Unix(int64(*blockchainCrowdloan.CloseDate), 0)
		loan.CloseDate = &closeDate
	}

	return loan
}

func CreateCrowdloanStatus(
	crowdloanUuid string,
	status string,
	date time.Time,
	comment string,
) (CrowdloanStatus, error) {
	crowdloanStatus := CrowdloanStatus{
		CrowdloanUuid: crowdloanUuid,
		Status:        status,
		Time:          date,
		Comment:       comment,
	}
	return crowdloanStatus, crowdloanStatus.Save()
}

func FindAllCrowdloans() Crowdloans {
	var crowdloans Crowdloans
	database.
		Preload("User").
		Preload("Status").
		Order("created_at DESC").
		Find(&crowdloans)

	return crowdloans
}

/*
	Blockchain Synchronization
*/

func (crowdloan *Crowdloan) UpdateStatus() error {
	if crowdloan.Type == "polkadot" {

		crowdloanInfo, err := apis.CrowdloanInfo(crowdloan.PolkadotCrowdloanId)
		if err != nil {
			return err
		}

		if crowdloan.Value != crowdloanInfo.Value {
			crowdloan.Value = crowdloanInfo.Value
			crowdloan.Save()
		}

		if crowdloan.ClosingAmount == nil || *crowdloan.ClosingAmount != crowdloanInfo.ClosingAmount {
			crowdloan.ClosingAmount = &crowdloanInfo.ClosingAmount
			crowdloan.Save()
		}

		if len(crowdloan.Status) == 0 ||
			crowdloan.Status[len(crowdloan.Status)-1].Status != crowdloanInfo.EscrowStatus.Status ||
			uint64(crowdloan.Status[len(crowdloan.Status)-1].Time.Unix()) != crowdloanInfo.EscrowStatus.Time {
			date := time.Unix(int64(crowdloanInfo.EscrowStatus.Time), 0)
			status, err := CreateCrowdloanStatus(
				crowdloan.Uuid,
				crowdloanInfo.EscrowStatus.Status,
				date,
				"updated from blockchain",
			)
			if err != nil {
				return err
			}
			crowdloan.Status = append(crowdloan.Status, status)
		}
	}

	return nil
}
