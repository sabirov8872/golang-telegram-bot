package models

type Welcome struct {
	ModelID       string         `json:"model_id"`
	Name          string         `json:"name"`
	PhotoSHA      string         `json:"photo_sha"`
	DealerExist   []string       `json:"dealer_exist"`
	Modifications []Modification `json:"modifications"`
	PhotoSha666   string         `json:"photo_sha666"`
}

type Modification struct {
	ModificationID string       `json:"modification_id"`
	Name           string       `json:"name"`
	Producing      string       `json:"producing"`
	Price          string       `json:"price"`
	Options        []string     `json:"options"`
	OptionsObj     []OptionsObj `json:"options_obj"`
	Colors         []Color      `json:"colors"`
}

type Color struct {
	ColorID     string       `json:"color_id"`
	Name        string       `json:"name"`
	HexValue    string       `json:"hex_value"`
	QueueNo     string       `json:"queue_no"`
	ExpectDate  string       `json:"expect_date"`
	PhotoSHA    string       `json:"photo_sha"`
	StockData   []StockDatum `json:"stock_data"`
	PhotoSha666 string       `json:"photo_sha666"`
}

type StockDatum struct {
	RegionID string `json:"region_id"`
	Stock    string `json:"stock"`
}

type OptionsObj struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Imagesha    string `json:"imagesha"`
}

type Client struct {
	UserID    int64
	FirstName string
	Subscribe bool
}
