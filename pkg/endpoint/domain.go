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
		if _, err := db.DB.Exec("INSERT INTO user (uid, guid, username, nickname, domain, passwd, role, authsource) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", util.GenUid(), state["global_login"], "root", "root", req.Domain.Uid, "", "ROOT", "WOJ"); err != nil {
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
	if state["global_login"] == "" {
		res.Success = false
		res.Error = "Login required"
		return
	}

	row := db.DB.QueryRowContext(ctx, "SELECT uid, description, domain, password, validto, invite_role, invite_state FROM domain_invite WHERE uid = ?", req.UID)

	ret := &base.DomainInvite{}

	if err := row.Scan(&ret.Uid, &ret.Description, &ret.Domain, &ret.Password, &ret.ValidTo, &ret.InviteRole, &ret.InviteState); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	if ret.Password != "" {
		ret.Password = "*"
	} else {
		ret.Password = ""
	}

	res.Success = true
	res.DomainInvite = ret
}

func getDomainInvites(ctx context.Context, req *senrenrpc.GetDomainInvitesRequest, state map[string]string, res *senrenrpc.GetDomainInvitesResponse) {
	if state["role"] != "ROOT" && state["role"] != "ADMIN" {
		res.Success = false
		res.Error = "forbidden"
		return
	}

	row, err := db.DB.QueryContext(ctx, "SELECT uid, description, domain, password, validto, invite_role, invite_state FROM domain_invite WHERE domain = ?", req.Domain)

	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	ret := make([]*base.DomainInvite, 0)

	for row.Next() {
		t := &base.DomainInvite{}
		row.Scan(&t.Uid, &t.Description, &t.Domain, &t.Password, &t.ValidTo, &t.InviteRole, &t.InviteState)
		ret = append(ret, t)
	}

	res.Success = true
	res.DomainInvites = ret
}

func createDomainInvite(ctx context.Context, req *senrenrpc.CreateDomainInviteRequest, state map[string]string, res *senrenrpc.CreateDomainInviteResponse) {
	if state["role"] != "ROOT" && state["role"] != "ADMIN" {
		res.Success = false
		res.Error = "forbidden"
		return
	}
	if (state["role"] == "ADMIN" && req.DomainInvite.InviteRole == "ADMIN") || req.DomainInvite.InviteRole == "ROOT" {
		res.Success = false
		res.Error = "you cannot create invitations with privilege higher than you"
		return
	}

	dbExec := "UPDATE domain_invite SET description = ? , password = ? , validto = ?, invite_role = ?, invite_state = ? WHERE uid = ? AND domain = ?"

	tDomain := senrenrpc.Domain(req.DomainInvite.Domain)
	tDomain.ConvertDomain()
	req.DomainInvite.Domain = string(tDomain)

	if req.DomainInvite.Uid == "" || req.DomainInvite.Uid == noUID {
		req.DomainInvite.Uid = util.GenUid()

		dbExec = "INSERT INTO domain_invite (description, password, validto, invite_role, invite_state, uid, domain) VALUES (?, ?, ?, ?, ?, ?, ?)"
	}

	_, err := db.DB.ExecContext(ctx, dbExec, req.DomainInvite.Description, req.DomainInvite.Password, req.DomainInvite.ValidTo, req.DomainInvite.InviteRole, req.DomainInvite.InviteState, req.DomainInvite.Uid, req.Domain.GetDomain())

	if err != nil {
		res.Error = err.Error()
		res.Success = false
		return
	}

	res.Success = true
	res.Domain = req.Domain
	res.UID = req.DomainInvite.Uid
}

func joinDomain(ctx context.Context, req *senrenrpc.JoinDomainRequest, state map[string]string, res *senrenrpc.JoinDomainResponse) {
	if state["global_login"] == "" {
		res.Success = false
		res.Error = "Login required"
		return
	}

	row := db.DB.QueryRowContext(ctx, "SELECT uid, description, domain, validto, password, invite_role, invite_state FROM domain_invite WHERE uid = ? AND domain = ?", req.InviteCode, req.Domain)

	ret := &base.DomainInvite{}

	if err := row.Scan(&ret.Uid, &ret.Description, &ret.Domain, &ret.ValidTo, &ret.Password, &ret.InviteRole, &ret.InviteState); err != nil {
		res.Success = false
		res.Error = err.Error()
		return
	}

	if req.InvitePassword != ret.Password {
		res.Success = false
		res.Error = "Wrong Password"
		return
	}

	_, err := db.DB.ExecContext(ctx, "INSERT INTO user (uid, guid, username, passwd, nickname, domain, role, state, authsource) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", util.GenUid(), state["global_login"], "."+state["gname"], "", req.NickName, ret.Domain, ret.InviteRole, ret.InviteState, "WOJ")
	if err != nil {
		res.Success = false
		res.Error = "Failed to join group, maybe you are already a group member (" + err.Error() + ")"
		return
	}

	res.Success = true
}

func getDomainUser(ctx context.Context, req *senrenrpc.GetDomainUserRequest, state map[string]string, res *senrenrpc.GetDomainUserResponse) {

}

func getDomainUsers(ctx context.Context, req *senrenrpc.GetDomainUsersRequest, state map[string]string, res *senrenrpc.GetDomainUsersResponse) {

}

func updateDomainUser(ctx context.Context, req *senrenrpc.UpdateDomainUserRequest, state map[string]string, res *senrenrpc.UpdateDomainUserResponse) {

}

func getPCISid(ctx context.Context, req *senrenrpc.GetPCISidRequest, state map[string]string, res *senrenrpc.GetPCISidResponse) {
	if state["role"] == "" || state["role"] == "NONE" {
		res.Success = false
		res.Error = "not member of expected group"
		return
	}

	res.Session.Sid = util.SignSessionDomain(state["name"], req.GetDomain())
	res.Success = true
}
