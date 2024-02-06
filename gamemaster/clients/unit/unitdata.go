package unit

type Position struct {
	Col int `json:"Col"`
	Row int `json:"Row"`
}
type Location struct {
	BattlefieldID string   `json:"battlefieldID"`
	Position      Position `json:"position"`
	Heading       int      `json:"heading"`
}
type Stats struct {
	Name       string   `json:"name"`
	Model      string   `json:"model"`
	Pointvalue int      `json:"pointvalue"`
	Type       string   `json:"type"`
	Size       int      `json:"size"`
	Movement   string   `json:"movement"`
	Shortdmg   int      `json:"shortdmg"`
	Meddmg     int      `json:"meddmg"`
	Longdmg    int      `json:"longdmg"`
	Ovhdmg     int      `json:"ovhdmg"`
	Struct     int      `json:"struct"`
	Specials   []string `json:"specials"`
	Image      string   `json:"image"`
}
type UnitData struct {
	Location Location `json:"location"`
	Stats    Stats    `json:"stats"`
	Owner    string   `json:"owner"`
	Active   bool     `json:"active"`
}
