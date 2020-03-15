package style

type Table struct {
	ColumnWidth float64 `style:"column-width"`
	ColumnSpan  int     `style:"column-span"`
	RowSpan     int     `style:"row-span"`
}
