package krabtpl

import (
	"fmt"
	"time"

	"github.com/jaswdr/faker"
	"github.com/wzshiming/ctc"
)

// Fake generates fake data
func Fake(fake faker.Faker) any {
	return func(id string) string {
		switch id {
		case "Address.Country":
			return fake.Address().Country()
		case "Address.CountryCode":
			return fake.Address().CountryCode()

		case "Color.Name":
			return fake.Color().ColorName()
		case "Color.Hex":
			return fake.Color().Hex()
		case "Color.RGB":
			return fake.Color().RGB()

		case "Company.Name":
			return fake.Company().Name()

		case "Emoji.Emoji":
			return fake.Emoji().Emoji()
		case "Emoji.EmojiCode":
			return fake.Emoji().EmojiCode()

		case "File.Path":
			return fake.File().AbsoluteFilePathForUnix(4)
		case "File.Extension":
			return fake.File().Extension()
		case "File.FilenameWithExtension":
			return fake.File().FilenameWithExtension()
		case "File.MimeType":
			return fake.MimeType().MimeType()

		case "Food.Fruit":
			return fake.Food().Fruit()
		case "Food.Vegetable":
			return fake.Food().Vegetable()

		case "Hash.MD5":
			return fake.Hash().MD5()
		case "Hash.SHA256":
			return fake.Hash().SHA256()

		case "Internet.Domain":
			return fake.Internet().Domain()
		case "Internet.Email":
			return fake.Internet().Email()
		case "Internet.Ipv4":
			return fake.Internet().Ipv4()
		case "Internet.Ipv6":
			return fake.Internet().Ipv6()
		case "Internet.MacAddress":
			return fake.Internet().MacAddress()
		case "Internet.SafeEmail":
			return fake.Internet().SafeEmail()
		case "Internet.Slug":
			return fake.Internet().Slug()
		case "Internet.URL":
			return fake.Internet().URL()
		case "Internet.UserAgent":
			return fake.UserAgent().UserAgent()

		case "Lorem.Paragraph":
			return fake.Lorem().Paragraph(5)
		case "Lorem.Word":
			return fake.Lorem().Word()

		case "Person.FirstName":
			return fake.Person().FirstName()
		case "Person.LastName":
			return fake.Person().LastName()
		case "Person.Name":
			return fake.Person().Name()

		case "Time.ISO8601":
			return fake.Time().ISO8601(time.Now())
		case "Time.Timezone":
			return fake.Time().Timezone()
		}

		panic(fmt.Sprintf("Invalid generator name: %s%s%s", ctc.ForegroundRed, id, ctc.Reset))
	}
}
