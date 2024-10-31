package clients

import (
	"context"
	"fmt"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/rcdevice"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/token"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	c *resty.Client
}

func NewClient(address string) *Client {
	c := resty.New()
	c.BaseURL = address
	c.Header.Add("Content-Type", "application/json")
	c.OnAfterResponse(
		func(c *resty.Client, r *resty.Response) error {
			if r.StatusCode()/100 != 2 {
				return fmt.Errorf("resp not 2xx,%s", string(r.Body()))
			}
			return nil
		},
	)
	return &Client{
		c: c,
	}
}

func (c *Client) Debug(v bool) {
	c.c.Debug = v
}

func (c *Client) Auth(username, password string) error {
	tk := token.DefaultToken()
	_, err := c.c.R().SetContext(context.Background()).SetBody(token.NewIssueTokenRequest(username, password)).SetResult(tk).Post("/wod/api/v1/tokens")
	if err != nil {
		return err
	}
	c.c.SetAuthToken(tk.AccessToken)
	return nil
}

func (c *Client) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	us := user.NewUser(req)
	_, err := c.c.R().SetContext(ctx).SetResult(us).Post("/wod/api/v1/users")
	if err != nil {
		return nil, err
	}
	return us, nil
}

func (c *Client) ChangeUser(ctx context.Context, req *user.ChangeUserRequest) (*user.User, error) {
	nusreq := user.NewCreateUserRuquest()
	us := user.NewUser(nusreq)
	_, err := c.c.R().SetContext(ctx).SetResult(us).Patch("/wod/api/v1/users")
	if err != nil {
		return nil, err
	}
	return us, nil
}

func (c *Client) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.User, error) {
	nusreq := user.NewCreateUserRuquest()
	us := user.NewUser(nusreq)
	_, err := c.c.R().SetContext(ctx).SetResult(us).Delete("/wod/api/v1/users")
	if err != nil {
		return nil, err
	}
	return us, nil
}

func (c *Client) QueryDeviceList(ctx context.Context, req *rcdevice.QueryDeviceListRequest) (*rcdevice.DeviceSet, error) {
	ds := rcdevice.NewDeviceSet()
	_, err := c.c.R().SetContext(ctx).SetResult(ds).Get("/wod/api/v1/rcdevice")
	if err != nil {
		return nil, err
	}
	return ds, nil
}
