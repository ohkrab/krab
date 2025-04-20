package tpls

import (
	"fmt"
	"time"

	"github.com/jaswdr/faker"
	"github.com/wzshiming/ctc"
)

type Faker struct {
	fake faker.Faker
}

func NewFaker() *Faker {
	return &Faker{fake: faker.New()}
}

// Get generates fake data
func (f *Faker) Get() func(id string) string {
	return func(id string) string {
		switch id {
		case "Address.Country":
			return f.fake.Address().Country()
		case "Address.CountryCode":
			return f.fake.Address().CountryCode()

		case "Color.Name":
			return f.fake.Color().ColorName()
		case "Color.Hex":
			return f.fake.Color().Hex()
		case "Color.RGB":
			return f.fake.Color().RGB()

		case "Company.Name":
			return f.fake.Company().Name()

		case "Emoji.Emoji":
			return f.fake.Emoji().Emoji()
		case "Emoji.EmojiCode":
			return f.fake.Emoji().EmojiCode()

		case "File.Path":
			return f.fake.File().AbsoluteFilePathForUnix(4)
		case "File.Extension":
			return f.fake.File().Extension()
		case "File.FilenameWithExtension":
			return f.fake.File().FilenameWithExtension()
		case "File.MimeType":
			return f.fake.MimeType().MimeType()

		case "Food.Fruit":
			return f.fake.Food().Fruit()
		case "Food.Vegetable":
			return f.fake.Food().Vegetable()

		case "Hash.MD5":
			return f.fake.Hash().MD5()
		case "Hash.SHA256":
			return f.fake.Hash().SHA256()

		case "Internet.Domain":
			return f.fake.Internet().Domain()
		case "Internet.Email":
			return f.fake.Internet().Email()
		case "Internet.Ipv4":
			return f.fake.Internet().Ipv4()
		case "Internet.Ipv6":
			return f.fake.Internet().Ipv6()
		case "Internet.MacAddress":
			return f.fake.Internet().MacAddress()
		case "Internet.SafeEmail":
			return f.fake.Internet().SafeEmail()
		case "Internet.Slug":
			return f.fake.Internet().Slug()
		case "Internet.URL":
			return f.fake.Internet().URL()
		case "Internet.UserAgent":
			return f.fake.UserAgent().UserAgent()

		case "Lorem.Paragraph":
			return f.fake.Lorem().Paragraph(5)
		case "Lorem.Word":
			return f.fake.Lorem().Word()

		case "Person.FirstName":
			return f.fake.Person().FirstName()
		case "Person.LastName":
			return f.fake.Person().LastName()
		case "Person.Name":
			return f.fake.Person().Name()

		case "Time.ISO8601":
			return f.fake.Time().ISO8601(time.Now())
		case "Time.Timezone":
			return f.fake.Time().Timezone()
		}

		panic(fmt.Sprintf("Invalid generator name: %s%s%s", ctc.ForegroundRed, id, ctc.Reset))
	}
}
