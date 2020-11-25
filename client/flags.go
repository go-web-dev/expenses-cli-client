package client

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// titleFlag represents the title flag
type titleFlag struct {
	value    string
	optional bool
}

func (t titleFlag) String() string {
	return t.value
}

func (t *titleFlag) Set(title string) error {
	title = strings.TrimSpace(title)
	if t.optional && title == "" {
		return nil
	}
	if title == "" {
		return errors.New("expense title flag is required and must not be empty")
	}
	t.value = title
	return nil
}

// setTitleFlag configures the title flag on a specific command
func setTitleFlag(f *flag.FlagSet, optional bool) *titleFlag {
	t := titleFlag{
		optional: optional,
	}
	description := "Expense title"
	f.Var(&t, "title", description)
	f.Var(&t, "t", description)
	return &t
}

// currencyFlag represents the currency flag
type currencyFlag struct {
	value    string
	optional bool
}

func (c currencyFlag) String() string {
	return c.value
}

func (c *currencyFlag) Set(currency string) error {
	currencies := []string{"USD", "EUR", "GBP", "MDL"}
	currency = strings.TrimSpace(strings.ToUpper(currency))
	if c.optional && currency == "" {
		return nil
	}
	switch currency {
	case currencies[0], currencies[1], currencies[2], currencies[3]:
		c.value = currency
		return nil
	default:
		return errors.New("currency must be one of: " + strings.Join(currencies, ","))
	}
}

// setCurrencyFlag configures the currency flag for a specific command
func setCurrencyFlag(f *flag.FlagSet, optional bool) *currencyFlag {
	c := currencyFlag{
		optional: optional,
	}
	description := "Expense currency"
	f.Var(&c, "currency", description)
	f.Var(&c, "c", description)
	return &c
}

// priceFlag represents the price flag
type priceFlag struct {
	value    float64
	optional bool
}

func (p priceFlag) String() string {
	return fmt.Sprintf("%f", p.value)
}

func (p *priceFlag) Set(v string) error {
	if p.optional && v == "" || v == "0" {
		return nil
	}
	price, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err
	}
	if price <= 0 {
		return errors.New("expense price is required and must be bigger than 0")
	}
	p.value = price
	return nil
}

// setCurrencyFlag configures the currency flag for a specific command
func setPriceFlag(f *flag.FlagSet, optional bool) *priceFlag {
	p := priceFlag{
		optional: optional,
	}
	description := "Expense price"
	f.Var(&p, "price", description)
	f.Var(&p, "p", description)
	return &p
}

// ids represents a list of ids (strings)
type idsFlag struct {
	value  []string
	unique map[string]struct{}
}

func (ids idsFlag) String() string {
	return strings.Join(ids.value, ",")
}

func (ids *idsFlag) Set(id string) error {
	uid, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return errors.Wrap(err, "invalid uuid provided")
	}
	_, found := ids.unique[uid.String()]
	if !found {
		ids.value = append(ids.value, uid.String())
		ids.unique[uid.String()] = struct{}{}
	}
	return nil
}

// setIDsFlag configures the ids flag for a specific command
func setIDsFlag(f *flag.FlagSet) *idsFlag {
	ids := idsFlag{
		unique: map[string]struct{}{},
	}
	description := "Expense id o be appended to the list of IDs"
	f.Var(&ids, "id", description)
	return &ids
}

// pageFlag represents the page flag
type pageFlag struct {
	value string
}

func (p pageFlag) String() string {
	return p.value
}

func (p *pageFlag) Set(page string) error {
	page = strings.TrimSpace(page)
	if page == "" {
		p.value = "1"
		return nil
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return errors.Wrap(err, "invalid page provided")
	}
	if pageInt <= 0 {
		return errors.New("page must be greater than 0")
	}
	p.value = page
	return nil
}

// setPageFlag configures the page flag on a specific command
func setPageFlag(f *flag.FlagSet) *pageFlag {
	var p pageFlag
	description := "Page number (for pagination)"
	f.Var(&p, "page", description)
	f.Var(&p, "p", description)
	return &p
}

// pageSizeFlag represents the page_size flag
type pageSizeFlag struct {
	value string
}

func (ps pageSizeFlag) String() string {
	return ps.value
}

func (ps *pageSizeFlag) Set(pageSize string) error {
	pageSize = strings.TrimSpace(pageSize)
	if pageSize == "" {
		ps.value = "5"
		return nil
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return errors.Wrap(err, "invalid page provided")
	}
	if pageSizeInt <= 0 || pageSizeInt > 25 {
		return errors.New("page_size must be greater than 0 and less than 25")
	}
	ps.value = pageSize
	return nil
}

// setPageSizeFlag configures the page_size flag on a specific command
func setPageSizeFlag(f *flag.FlagSet) *pageSizeFlag {
	var ps pageSizeFlag
	description := "Page size (for pagination)"
	f.Var(&ps, "page_size", description)
	f.Var(&ps, "ps", description)
	return &ps
}

// emailFlag represents the email flag
type emailFlag struct {
	value string
}

func (e emailFlag) String() string {
	return e.value
}

func (e *emailFlag) Set(email string) error {
	email = strings.TrimSpace(email)
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) < 3 && len(email) > 254 {
		return errors.New("the length of the email must be >= 3 and <= 254")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email address format")
	}
	e.value = email
	return nil
}

// setEmailFlag configures the email flag on a specific command
func setEmailFlag(f *flag.FlagSet) *emailFlag {
	var e emailFlag
	description := "The user email for login/signup"
	f.Var(&e, "email", description)
	f.Var(&e, "e", description)
	return &e
}

// passwordFlag represents the password flag
type passwordFlag struct {
	value string
}

func (e passwordFlag) String() string {
	return e.value
}

func (e *passwordFlag) Set(password string) error {
	if err := verifyPwd(password); err != nil {
		return err
	}
	e.value = password
	return nil
}

// setPasswordFlag configures the password flag on a specific command
func setPasswordFlag(f *flag.FlagSet) *passwordFlag {
	var p passwordFlag
	description := "The user password for login/signup"
	f.Var(&p, "password", description)
	f.Var(&p, "p", description)
	return &p
}

func verifyPwd(pwd string) error {
	var number, upper, lower, special bool
	for _, c := range pwd {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLower(c) || c == ' ':
			lower = true
		}
	}
	if !number {
		return errors.New("the password must have at least 1 digit")
	}
	if !upper {
		return errors.New("the password must have at least 1 uppercase letter")
	}
	if !lower {
		return errors.New("the password must have at least 1 lowercase letter")
	}
	if !special {
		return errors.New("the password must have at least 1 special symbol")
	}
	if len(pwd) < 8 && len(pwd) > 32 {
		return errors.New("the password length must be >= 8 && <= 32")
	}
	return nil
}
