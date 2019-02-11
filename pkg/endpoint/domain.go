package endpoint

import (
	"context"

	"github.com/erjiaqing/senren2/pkg/db"
	"github.com/erjiaqing/senren2/pkg/types/base"
	"github.com/erjiaqing/senren2/pkg/types/senrenrpc"
	"github.com/erjiaqing/senren2/pkg/util"
)

func getDomain(ctx context.Context, req *senrenrpc.GetDomainRequest, state map[string]string, res *senrenrpc.GetDomainResponse) {
	ret := &base.DomainInfo{}

	if req.Domain == "aaaaaaaaaaaaaaaa" {
		res.Success = true
		res.Domain = ret
		return
	}

	row := db.DB.QueryRowContext(ctx, "SELECT title, is_public, uid, short_name, description FROM domain_info WHERE uid = ?", state["domain"])
	if err := row.Scan(&ret.Title, &ret.IsPublic, &ret.Uid, &ret.ShortName, &ret.Description); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	if req.Filter == "nodesc" {
		ret.Description = ""
	}
	res.Success = true
	res.Domain = ret
}

func getDomains(ctx context.Context, req *senrenrpc.GetDomainsRequest, state map[string]string, res *senrenrpc.GetDomainsResponse) {
	ret := make([]*base.DomainInfo, 0)
	rows, err := db.DB.QueryContext(ctx, "SELECT uid, short_name, title FROM domain_info WHERE is_public = ?", "PUBLIC")
	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	for rows.Next() {
		t := &base.DomainInfo{}
		rows.Scan(&t.Uid, &t.ShortName, &t.Title)
		ret = append(ret, t)
	}
	res.Success = true
	res.Domains = ret
}

func createDomain(ctx context.Context, req *senrenrpc.CreateDomainRequest, state map[string]string, res *senrenrpc.CreateDomainResponse) {
	dbExec := "UPDATE domain_info SET title = ? , description = ?, short_name = ?, is_public = ? WHERE uid = ?"
	doDomainCreate := false
	if req.Domain.Uid == "" || req.Domain.Uid == noUID {
		req.Domain.Uid = util.GenUid()
		dbExec = "INSERT INTO domain_info (title, description, short_name, is_public, uid) VALUES (?, ?, ?, ?, ?)"
		doDomainCreate = true
	}

	tDomain := senrenrpc.Domain(req.Domain.Uid)
	tDomain.ConvertDomain()
	req.Domain.Uid = tDomain.GetDomain()

	if req.Domain.IsPublic == "" {
		req.Domain.IsPublic = "PRIVATE"
		// Private is default, only VIP user can create Public and Protected groups
	}

	if req.Domain.IsPublic != "PUBLIC" && req.Domain.IsPublic != "PROTECTED" {
		req.Domain.ShortName = ""
		// non-public domain cannot have short name
		// PROTECTED: content is public but only domain members can interact (submit, discuess, etc.)
	}

	if doDomainCreate {
		if _, err := db.DB.Exec("INSERT INTO user (uid, guid, username, nickname, domain, password, role, authsource) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", util.GenUid(), state["global_login"], "root", "root", req.Domain.Uid, "", "ROOT", "WOJ"); err != nil {
			panic(err.Error())
		}
	}

	if _, err := db.DB.Exec(dbExec, req.Domain.Title, req.Domain.Description, req.Domain.ShortName, req.Domain.IsPublic, req.Domain.Uid); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}
	res.Domain = senrenrpc.Domain(req.Domain.Uid)
	res.UID = req.Domain.Uid
	res.Success = true
}

func getDomainInvite(ctx context.Context, req *senrenrpc.GetDomainInviteRequest, state map[string]string, res *senrenrpc.GetDomainInviteResponse) {

}

func getDomainInvites(ctx context.Context, req *senrenrpc.GetDomainInvitesRequest, state map[string]string, res *senrenrpc.GetDomainInvitesResponse) {

}

func createDomainInvite(ctx context.Context, req *senrenrpc.CreateDomainInviteRequest, state map[string]string, res *senrenrpc.CreateDomainInviteResponse) {

}

func joinDomain(ctx context.Context, req *senrenrpc.JoinDomainRequest, state map[string]string, res *senrenrpc.JoinDomainResponse) {

}

func getDomainUser(ctx context.Context, req *senrenrpc.GetDomainUserRequest, state map[string]string, res *senrenrpc.GetDomainUserResponse) {

}

func getDomainUsers(ctx context.Context, req *senrenrpc.GetDomainUsersRequest, state map[string]string, res *senrenrpc.GetDomainUsersResponse) {

}

func updateDomainUser(ctx context.Context, req *senrenrpc.UpdateDomainUserRequest, state map[string]string, res *senrenrpc.UpdateDomainUserResponse) {

}
