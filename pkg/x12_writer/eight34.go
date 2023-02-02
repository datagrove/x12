package x12_writer

type Eight34 struct {
	Index            string
	Bgn02ReferenceId string //
	Ref38Policy      string
	Sponsor          Name
	Payer            Name
	Insured          []Insured
}

func New834() *Eight34 {
	return &Eight34{
		Bgn02ReferenceId: "1",
	}
}

const (
	Verify  = "4"
	Change  = "2"
	Replace = "RX"
)

func (x *Eight34) Write(w *EdiWriter) {
	w.BeginGroup()
	defer w.EndGroup()
	w.BeginTransaction("834")
	defer w.EndTransaction()
	w.Write("BGN", "00", x.Bgn02ReferenceId, w.ccyymmdd, w.hhmm, "", "", "", Verify)
	w.Ref("38", x.Ref38Policy)
	w.Date("303", w.ccyymmdd)
	//w.Write("QTY", "TO", fmt.Sprintf("%d", len(x.Insured))
	x.Sponsor.N1("P5", w)
	x.Payer.N1("IN", w)
	for _, o := range x.Insured {
		o.Write(w)
	}
}

type Insured struct {
	Relationship string
	Name
	Address
	Policy    string
	Effective string
	Terminate string
	Demographic
	Phone
}

func (x *Insured) Write(w *EdiWriter) {
	isSubscriber := "Y"
	if x.Relationship != "18" {
		isSubscriber = "N"
	}

	w.Write("INS", isSubscriber, x.Relationship, "030", "20", "A")
	w.Write("REF", "0F", x.Id)
	x.Name.Nm1("IL", w)
	x.Phone.Per("IP", w)
	x.Address.N3N4(w)
	x.Demographic.Dmg(w)
	w.Write("HD", "030", "", "HLT", x.Policy, "EMP")
	w.Date("348", x.Effective)
	w.Date("349", x.Terminate)
	w.Ref("1L", x.Policy)
}
