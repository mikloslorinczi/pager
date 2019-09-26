package rules

var workHours = map[int]bool{
	0:  false,
	1:  false,
	2:  false,
	3:  false,
	4:  false,
	5:  false,
	6:  false,
	7:  false,
	8:  true,
	9:  true,
	10: true,
	11: true,
	12: true,
	13: true,
	14: true,
	15: true,
	16: true,
	17: false,
	18: false,
	19: false,
	20: false,
	21: false,
	22: false,
	23: false,
}

type dayException struct {
	year      int
	month     int
	day       int
	isWorkDay bool
}

var dayExceptions = []dayException{
	// 2019
	{
		year:      2019,
		month:     1,
		day:       1,
		isWorkDay: false,
	},
	{
		year:      2019,
		month:     3,
		day:       15,
		isWorkDay: false,
	},
	{
		year:      2019,
		month:     4,
		day:       19,
		isWorkDay: false,
	},
	{
		year:      2019,
		month:     4,
		day:       22,
		isWorkDay: false,
	},
	{
		year:      2019,
		month:     5,
		day:       1,
		isWorkDay: false,
	},
	{
		year:      2019,
		month:     6,
		day:       10,
		isWorkDay: false,
	},
	{
		year:      2019,
		month:     8,
		day:       10,
		isWorkDay: true,
	},
	{
		year:      2019,
		month:     8,
		day:       19,
		isWorkDay: false,
	},
	{
		year:      2019,
		month:     8,
		day:       20,
		isWorkDay: false,
	},
	{
		year:      2019,
		month:     10,
		day:       23,
		isWorkDay: false,
	},
	{
		year:      2019,
		month:     11,
		day:       1,
		isWorkDay: false,
	},
	{
		year:      2019,
		month:     12,
		day:       7,
		isWorkDay: true,
	},
	{
		year:      2019,
		month:     12,
		day:       14,
		isWorkDay: true,
	},
	{
		year:      2019,
		month:     12,
		day:       24,
		isWorkDay: false,
	},
	{
		year:      2019,
		month:     12,
		day:       25,
		isWorkDay: false,
	},
	{
		year:      2019,
		month:     12,
		day:       26,
		isWorkDay: false,
	},
	{
		year:      2019,
		month:     12,
		day:       27,
		isWorkDay: false,
	},
	// 2020
	{
		year:      2020,
		month:     1,
		day:       1,
		isWorkDay: false,
	},
	{
		year:      2020,
		month:     4,
		day:       10,
		isWorkDay: false,
	},
	{
		year:      2020,
		month:     4,
		day:       13,
		isWorkDay: false,
	},
	{
		year:      2020,
		month:     5,
		day:       1,
		isWorkDay: false,
	},
	{
		year:      2020,
		month:     6,
		day:       1,
		isWorkDay: false,
	},
	{
		year:      2020,
		month:     8,
		day:       20,
		isWorkDay: false,
	},
	{
		year:      2020,
		month:     8,
		day:       21,
		isWorkDay: false,
	},
	{
		year:      2020,
		month:     8,
		day:       29,
		isWorkDay: true,
	},
	{
		year:      2020,
		month:     10,
		day:       23,
		isWorkDay: false,
	},
	{
		year:      2020,
		month:     12,
		day:       12,
		isWorkDay: true,
	},
	{
		year:      2020,
		month:     12,
		day:       24,
		isWorkDay: false,
	},
	{
		year:      2020,
		month:     12,
		day:       25,
		isWorkDay: false,
	},
}
