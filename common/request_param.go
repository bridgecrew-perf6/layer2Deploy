package common

type TakeOrderParam struct {
	ProductName      string `json:"productName"`
	OntId            string `json:"ontId"`
	UserName         string `json:"userName"`
	ApiId            uint32 `json:"apiId"`
	SpecificationsId uint32 `json:"specificationsId"`
}

type RenewOrderParam struct {
	OrderId          string `json:"orderId"`
	SpecificationsId uint32 `json:"specificationsId"`
}

type AliPayOderParam struct {
	OrderId string `json:"orderId"`
}

type GenerateTestKeyParam struct {
	ApiId uint32 `json:"apiId"`
}

type AdminGenerateTestKeyParam struct {
	ApiId    uint32 `json:"apiId"`
	ApiState int32  `json:"apiState"`
}

type GetPayQrCode struct {
	OrderId  string `json:"orderId"`
	OntId    string `json:"ontId"`
	UserName string `json:"userName"`
}

type SendTxParam struct {
	Signer    string    `json:"signer"`
	SignedTx  string    `json:"signedTx"`
	ExtraData ExtraData `json:"extraData"`
}

type ExtraData struct {
	Id        string `json:"id"`
	PublicKey string `json:"publickey"`
	OntId     string `json:"ontId"`
}

type GetQrCodeParam struct {
	Id string `json:"id"`
}

type OrderIdParam struct {
	OrderId string `json:"orderId"`
}
type QrCodeIdParam struct {
	QrCodeId string `json:"qrCodeId"`
}

type SearchApiByKey struct {
	Key string `json:"key"`
}

type GetApiByCategoryId struct {
	CategoryId uint32 `json:"categoryId"`
	PageNum    uint32 `json:"pageNum"`
	PageSize   uint32 `json:"pageSize"`
}

type TokenTemplate struct {
	DataIds   string `json:"dataIds"`
	TokenHash string `json:"tokenHash"`
}

type Dtoken struct {
	ResourceId    string        `json:"resourceId"`
	Account       string        `json:"account"`
	TokenTemplate TokenTemplate `json:"tokenTemplate"`
	Number        uint32        `json:"number"`
}

type WetherForcastServiceRequest struct {
	DataUrl       string                 `json:"dataUrl"`
	Header        map[string]interface{} `json:"header"`
	Param         map[string]interface{} `json:"param"`
	RequestMethod string                 `json:"requestMethod"`
	AlgorithmName string                 `json:"algorithmName"`
	TokenId       string                 `json:"tokenId"`
	OwnerAddress  string                 `json:"ownerAddress"`
}

type OrderListRequest struct {
	Title    string `json:"title"`
	ApiId    uint32 `json:"apiId"`
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Sorting  string `json:"sorting"`
}

type CommentOrderApiRequest struct {
	OrderId  string `json:"orderId"`
	StarNum  uint32 `json:"starNum"`
	Comments string `json:"comment"`
}

type DelCommentOrderRequest struct {
	Id uint32 `json:"id"`
}
