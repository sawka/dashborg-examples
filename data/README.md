# Data Files


## yc-companies

Schema for yc-companies.db (sqlite3 db):

```
CREATE TABLE companies (
    id integer not null primary key autoincrement,
    name text not null,
    companyurl text,
    logourl text,
    rank int,
    overview text,
    sector text,
    jobscreated int,
    batch text,
    hqlocation text,
    jobsurl text
);
```

Golang struct (can read/write json or SQL):
```
type Company struct {
    Id          int `json:"id"`
	Name        string `json:"name"`
	CompanyUrl  string `json:"companyurl"`
	LogoUrl     string `json:"logourl"`
	Rank        int    `json:"rank"`
	Overview    string `json:"overview"`
	Sector      string `json:"sector"`
	JobsCreated int    `json:"jobscreated"`
	Batch       string `json:"batch"`
	HQLocation  string `json:"hqlocation"`
	JobsUrl     string `json:"jobsurl"`
}
```

Example entry from yc-companies.json:
```
{
  "name":"Airbnb",
  "companyurl":"http://airbnb.com",
  "logourl":"https://bookface-images.s3.amazonaws.com/small_logos/3e9a0092bee2ccf926e650e59c06503ec6b9ee65.png",
  "rank":1,
  "overview":"Book accommodations around the world.",
  "sector":"Travel,Leisure and Tourism",
  "jobscreated":6000,
  "batch":"W09",
  "hqlocation":"San Francisco",
  "jobsurl":"https://careers.airbnb.com/"
}
```
