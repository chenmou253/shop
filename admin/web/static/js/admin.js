(function () {
  function on(selector, event, handler) {
    document.addEventListener(event, function (e) {
      var target = e.target.closest(selector);
      if (target) handler(e, target);
    });
  }

  function api(path, options) {
    options = options || {};
    var headers = options.headers || {};
    if (options.body && !(options.body instanceof FormData)) {
      headers["Content-Type"] = "application/json";
    }
    return fetch(path, Object.assign({ headers: headers }, options)).then(function (resp) {
      if (resp.status === 401 && !document.querySelector("[data-login-form]")) {
        window.location.href = "/login";
        return Promise.reject(new Error("未登录"));
      }
      return resp.json().catch(function () { return {}; }).then(function (payload) {
        if (!resp.ok) throw new Error(payload.error || "请求失败");
        return payload;
      });
    });
  }

  function value(obj, lower, upper, fallback) {
    if (!obj) return fallback == null ? "" : fallback;
    if (obj[lower] != null) return obj[lower];
    if (obj[upper] != null) return obj[upper];
    return fallback == null ? "" : fallback;
  }

  function text(value) {
    return String(value == null ? "" : value)
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;");
  }

  function statusBadge(enabled) {
    return '<span class="badge ' + (enabled ? "success" : "secondary") + '">' + (enabled ? "启用" : "禁用") + "</span>";
  }

  function methodBadge(method) {
    var cls = { GET: "success", POST: "primary", PUT: "warning", PATCH: "warning", DELETE: "danger" }[method] || "secondary";
    return '<span class="badge ' + cls + '">' + text(method || "-") + "</span>";
  }

  function checkedValues(selector) {
    return Array.prototype.slice.call(document.querySelectorAll(selector))
      .filter(function (input) { return input.checked; })
      .map(function (input) { return Number(input.value); })
      .filter(Boolean);
  }

  function setAlert(message, type) {
    if (!message) return;
    var content = document.querySelector(".content") || document.body;
    var old = content.querySelector(".runtime-alert");
    if (old) old.remove();
    var div = document.createElement("div");
    div.className = "alert runtime-alert " + (type || "success");
    div.textContent = message;
    content.prepend(div);
    setTimeout(function () { div.remove(); }, 2800);
  }

  function fillForm(form, data) {
    Object.keys(data || {}).forEach(function (key) {
      var input = form.elements[key];
      if (!input) return;
      if (input.type === "checkbox") {
        input.checked = !!data[key];
      } else {
        input.value = data[key] == null ? "" : data[key];
      }
    });
  }

  function formBool(form, name) {
    return !!form.elements[name] && form.elements[name].checked;
  }

  function formNumber(form, name) {
    var n = Number(form.elements[name] ? form.elements[name].value : 0);
    return Number.isFinite(n) ? n : 0;
  }

  function flattenTree(nodes, pick) {
    var out = [];
    function walk(items, level) {
      (items || []).forEach(function (node) {
        out.push({ item: pick(node), level: level, node: node });
        walk(node.children || [], level + 1);
      });
    }
    walk(nodes || [], 0);
    return out;
  }

  function indent(level) {
    return new Array(level + 1).join("&nbsp;&nbsp;&nbsp;&nbsp;");
  }

  function initShell() {
    if (!document.body.classList.contains("admin-body")) return;
    api("/api/me").then(function (payload) {
      var user = payload.user || {};
      var target = document.querySelector("[data-current-user]");
      if (target) target.textContent = user.username || "";
    });
    api("/api/nav").then(function (payload) {
      renderNav(payload.menus || []);
    });
  }

  function renderNav(nodes) {
    var nav = document.querySelector("#sidebar-nav");
    if (!nav) return;
    var active = document.body.getAttribute("data-active") || "";
    var html = '<a class="nav-item ' + (active === "dashboard" ? "active" : "") + '" href="/dashboard" data-nav-key="dashboard"><span class="nav-icon">▣</span><span>控制台</span></a>';
    html += renderNavNodes(nodes, active);
    nav.innerHTML = html;
  }

  function renderNavNodes(nodes, active) {
    return (nodes || []).map(function (node) {
      var menu = node.menu || {};
      var title = value(menu, "title", "Title");
      var icon = value(menu, "icon", "Icon", "▦");
      var path = value(menu, "path", "Path", "#");
      var children = node.children || [];
      if (children.length) {
        return '<div class="nav-branch open"><button class="nav-item nav-toggle" type="button" data-tree-toggle>' +
          '<span class="nav-icon">' + text(icon) + '</span><span>' + text(title) + '</span><span class="nav-caret">▾</span></button>' +
          '<div class="nav-children">' + renderNavNodes(children, active) + "</div></div>";
      }
      var key = path.replace("/", "") || "dashboard";
      return '<a class="nav-item ' + (active === key ? "active" : "") + '" href="' + text(path) + '">' +
        '<span class="nav-icon">' + text(icon || "•") + '</span><span>' + text(title) + "</span></a>";
    }).join("");
  }

  function initLogin() {
    var form = document.querySelector("[data-login-form]");
    if (!form) return;
    form.addEventListener("submit", function (e) {
      e.preventDefault();
      api("/api/login", {
        method: "POST",
        body: JSON.stringify({
          username: form.elements.username.value,
          password: form.elements.password.value
        })
      }).then(function () {
        window.location.href = "/dashboard";
      }).catch(function (err) {
        setAlert(err.message, "danger");
      });
    });
  }

  function initDashboard() {
    if (!document.querySelector("[data-dashboard]")) return;
    api("/api/dashboard").then(function (payload) {
      Object.keys(payload).forEach(function (key) {
        var target = document.querySelector('[data-stat="' + key + '"]');
        if (target) target.textContent = payload[key];
      });
    });
  }

  function initUsersIndex() {
    var body = document.querySelector("[data-users-table]");
    if (!body) return;
    loadUsers();
  }

  function loadUsers() {
    var body = document.querySelector("[data-users-table]");
    api("/api/users").then(function (payload) {
      var users = payload.users || [];
      body.innerHTML = users.length ? users.map(function (user) {
        var roles = (user.roles || []).map(function (role) { return role.name; }).join("、") || "-";
        return "<tr><td>" + user.id + "</td><td><strong>" + text(user.username) + "</strong></td>" +
          "<td>" + text(user.nickname) + "</td><td>" + text(user.email) + "</td><td>" + text(roles) + "</td>" +
          "<td>" + statusBadge(user.status) + '</td><td class="actions">' +
          '<a class="btn small" href="/users/' + user.id + '/edit">编辑</a>' +
          '<button class="btn small danger" type="button" data-delete-user="' + user.id + '">删除</button>' +
          "</td></tr>";
      }).join("") : '<tr><td colspan="7" class="empty-cell">暂无用户</td></tr>';
    }).catch(function (err) {
      body.innerHTML = '<tr><td colspan="7" class="empty-cell">' + text(err.message) + "</td></tr>";
    });
  }

  function initUserForm() {
    var form = document.querySelector("[data-user-form]");
    if (!form) return;
    var mode = form.getAttribute("data-mode");
    var id = form.getAttribute("data-id");
    var selected = [];
    Promise.all([
      api("/api/roles"),
      mode === "edit" ? api("/api/users/" + id) : Promise.resolve({ user: { status: true } })
    ]).then(function (items) {
      var roles = items[0].roles || [];
      var user = items[1].user || {};
      selected = user.role_ids || [];
      fillForm(form, user);
      renderUserRoles(roles, selected);
    }).catch(function (err) { setAlert(err.message, "danger"); });

    form.addEventListener("submit", function (e) {
      e.preventDefault();
      var payload = {
        username: form.elements.username.value.trim(),
        password: form.elements.password.value,
        nickname: form.elements.nickname.value.trim(),
        email: form.elements.email.value.trim(),
        status: formBool(form, "status"),
        role_ids: checkedValues('#user-role-options input[name="role_ids"]')
      };
      api(mode === "edit" ? "/api/users/" + id : "/api/users", {
        method: mode === "edit" ? "PUT" : "POST",
        body: JSON.stringify(payload)
      }).then(function () { window.location.href = "/users?flash=" + encodeURIComponent("保存成功"); })
        .catch(function (err) { setAlert(err.message, "danger"); });
    });
  }

  function renderUserRoles(roles, selected) {
    var box = document.querySelector("#user-role-options");
    if (!box) return;
    var chosen = new Set((selected || []).map(Number));
    box.innerHTML = roles.length ? roles.map(function (role) {
      return '<label class="check-item"><input type="checkbox" name="role_ids" value="' + role.id + '"' + (chosen.has(role.id) ? " checked" : "") + ">" +
        "<span>" + text(role.name) + "</span><em>" + text(role.code) + "</em></label>";
    }).join("") : '<div class="empty-cell">暂无可选角色</div>';
  }

  function initRolesIndex() {
    if (!document.querySelector("[data-roles-table]")) return;
    loadRoles();
  }

  function loadRoles() {
    var body = document.querySelector("[data-roles-table]");
    api("/api/roles").then(function (payload) {
      var roles = payload.roles || [];
      body.innerHTML = roles.length ? roles.map(function (role) {
        return "<tr><td>" + role.id + "</td><td><strong>" + text(role.name) + "</strong></td>" +
          "<td><code>" + text(role.code) + "</code></td><td>" + text(role.remark) + "</td><td>" + statusBadge(role.status) + '</td><td class="actions">' +
          '<a class="btn small" href="/roles/' + role.id + '/edit">编辑</a>' +
          '<button class="btn small danger" type="button" data-delete-role="' + role.id + '">删除</button>' +
          "</td></tr>";
      }).join("") : '<tr><td colspan="6" class="empty-cell">暂无角色</td></tr>';
    }).catch(function (err) {
      body.innerHTML = '<tr><td colspan="6" class="empty-cell">' + text(err.message) + "</td></tr>";
    });
  }

  function initRoleForm() {
    var form = document.querySelector("[data-role-form]");
    if (!form) return;
    var mode = form.getAttribute("data-mode");
    var id = form.getAttribute("data-id");
    var rolePromise = mode === "edit" ? api("/api/roles/" + id) : Promise.resolve({ role: { status: true, menu_ids: [], permission_ids: [] } });
    rolePromise.then(function (payload) {
      var role = payload.role || {};
      fillForm(form, role);
      return loadRoleAccess(id, role.menu_ids || [], role.permission_ids || []);
    }).catch(function (err) { setAlert(err.message, "danger"); });

    form.addEventListener("submit", function (e) {
      e.preventDefault();
      var payload = {
        name: form.elements.name.value.trim(),
        code: form.elements.code.value.trim(),
        remark: form.elements.remark.value.trim(),
        status: formBool(form, "status"),
        menu_ids: checkedValues('#role-menu-tree input[name="menu_ids"]'),
        permission_ids: collectPermissionIDs()
      };
      api(mode === "edit" ? "/api/roles/" + id : "/api/roles", {
        method: mode === "edit" ? "PUT" : "POST",
        body: JSON.stringify(payload)
      }).then(function () { window.location.href = "/roles?flash=" + encodeURIComponent("保存成功"); })
        .catch(function (err) { setAlert(err.message, "danger"); });
    });
  }

  function loadRoleAccess(roleID, selectedMenus, selectedPermissions) {
    var params = new URLSearchParams();
    params.set("role_id", roleID || "0");
    params.set("menu_id", document.querySelector("[data-role-access-filter]") ? document.querySelector("[data-role-access-filter]").value || "0" : "0");
    appendIDs(params, "selected_menus", selectedMenus || []);
    appendIDs(params, "selected_permissions", selectedPermissions || []);
    return api("/api/role-access?" + params.toString()).then(function (payload) {
      var menus = payload.menus || [];
      var perms = payload.permissions || [];
      renderRoleMenuFilter(menus);
      document.querySelector("#role-menu-tree").innerHTML = renderMenuChecks(menus);
      document.querySelector("#role-permission-tree").innerHTML = renderPermissionChecks(perms);
      syncHiddenSelected('#role-permission-tree input[name="permission_ids"]', "#role-hidden-permissions", "permission_ids", selectedPermissions || []);
    });
  }

  function renderRoleMenuFilter(nodes) {
    var select = document.querySelector("[data-role-access-filter]");
    if (!select || select.options.length > 1) return;
    flattenTree(nodes, function (node) { return node.menu || {}; }).forEach(function (row) {
      var menu = row.item;
      var option = document.createElement("option");
      option.value = value(menu, "id", "ID");
      option.innerHTML = indent(row.level) + text(value(menu, "title", "Title"));
      select.appendChild(option);
    });
  }

  function initMenusIndex() {
    if (!document.querySelector("[data-menus-table]")) return;
    loadMenus();
  }

  function loadMenus() {
    var body = document.querySelector("[data-menus-table]");
    api("/api/menus/tree").then(function (payload) {
      var rows = flattenTree(payload.tree || [], function (node) { return node.menu || {}; });
      body.innerHTML = rows.length ? rows.map(function (row) {
        var menu = row.item;
        var id = value(menu, "id", "ID");
        return '<tr data-level="' + row.level + '"><td>' + id + '</td><td>' + indent(row.level) + '<span class="tree-table-caret">▾</span><strong>' + text(value(menu, "title", "Title")) + "</strong></td>" +
          "<td><code>" + text(value(menu, "path", "Path")) + "</code></td><td>" + value(menu, "sort", "Sort", 0) + "</td>" +
          "<td>" + statusBadge(!!value(menu, "visible", "Visible")) + "</td><td>" + statusBadge(!!value(menu, "status", "Status")) + '</td><td class="actions">' +
          '<a class="btn small" href="/menus/' + id + '/edit">编辑</a>' +
          '<button class="btn small danger" type="button" data-delete-menu="' + id + '">删除</button>' +
          "</td></tr>";
      }).join("") : '<tr><td colspan="7" class="empty-cell">暂无菜单</td></tr>';
    }).catch(function (err) {
      body.innerHTML = '<tr><td colspan="7" class="empty-cell">' + text(err.message) + "</td></tr>";
    });
  }

  function initMenuForm() {
    var form = document.querySelector("[data-menu-form]");
    if (!form) return;
    var mode = form.getAttribute("data-mode");
    var id = form.getAttribute("data-id");
    Promise.all([
      api("/api/menus/tree"),
      mode === "edit" ? api("/api/menus/" + id) : Promise.resolve({ menu: { sort: 10, visible: true, status: true } })
    ]).then(function (items) {
      var menu = items[1].menu || {};
      renderParentMenuOptions(form.elements.parent_id, items[0].tree || [], menu.parent_id || 0, Number(id || 0));
      fillForm(form, menu);
    }).catch(function (err) { setAlert(err.message, "danger"); });
    form.addEventListener("submit", function (e) {
      e.preventDefault();
      var payload = {
        parent_id: formNumber(form, "parent_id"),
        title: form.elements.title.value.trim(),
        icon: form.elements.icon.value.trim(),
        path: form.elements.path.value.trim(),
        sort: formNumber(form, "sort"),
        visible: formBool(form, "visible"),
        status: formBool(form, "status")
      };
      api(mode === "edit" ? "/api/menus/" + id : "/api/menus", {
        method: mode === "edit" ? "PUT" : "POST",
        body: JSON.stringify(payload)
      }).then(function () { window.location.href = "/menus?flash=" + encodeURIComponent("保存成功"); })
        .catch(function (err) { setAlert(err.message, "danger"); });
    });
  }

  function renderParentMenuOptions(select, nodes, selected, currentID) {
    select.innerHTML = '<option value="0">顶级菜单</option>';
    flattenTree(nodes, function (node) { return node.menu || {}; }).forEach(function (row) {
      var menu = row.item;
      var id = Number(value(menu, "id", "ID"));
      if (id === currentID) return;
      var option = document.createElement("option");
      option.value = id;
      option.innerHTML = indent(row.level) + text(value(menu, "title", "Title"));
      option.selected = Number(selected) === id;
      select.appendChild(option);
    });
  }

  function initPermissionsIndex() {
    if (!document.querySelector("[data-permissions-table]")) return;
    loadPermissions();
  }

  function loadPermissions() {
    var body = document.querySelector("[data-permissions-table]");
    api("/api/permissions/tree").then(function (payload) {
      var rows = flattenTree(payload.tree || [], function (node) { return node.permission || {}; });
      body.innerHTML = rows.length ? rows.map(function (row) {
        var perm = row.item;
        var id = value(perm, "id", "ID");
        var menu = value(perm, "menu", "Menu", {});
        return '<tr data-level="' + row.level + '"><td>' + id + '</td><td>' + indent(row.level) + '<span class="tree-table-caret">▾</span><strong>' + text(value(perm, "name", "Name")) + "</strong></td>" +
          "<td><code>" + text(value(perm, "code", "Code")) + "</code></td><td>" + text(value(menu, "title", "Title")) + "</td>" +
          "<td>" + methodBadge(value(perm, "method", "Method")) + "</td><td><code>" + text(value(perm, "path", "Path")) + "</code></td>" +
          "<td>" + statusBadge(!!value(perm, "status", "Status")) + '</td><td class="actions">' +
          '<a class="btn small" href="/permissions/' + id + '/edit">编辑</a>' +
          '<button class="btn small danger" type="button" data-delete-permission="' + id + '">删除</button>' +
          "</td></tr>";
      }).join("") : '<tr><td colspan="8" class="empty-cell">暂无权限</td></tr>';
    }).catch(function (err) {
      body.innerHTML = '<tr><td colspan="8" class="empty-cell">' + text(err.message) + "</td></tr>";
    });
  }

  function initPermissionForm() {
    var form = document.querySelector("[data-permission-form]");
    if (!form) return;
    var mode = form.getAttribute("data-mode");
    var id = form.getAttribute("data-id");
    Promise.all([
      api("/api/menus/tree"),
      mode === "edit" ? api("/api/permissions/" + id) : Promise.resolve({ permission: { method: "GET", sort: 10, status: true, parent_id: 0, menu_id: 0 } })
    ]).then(function (items) {
      var perm = items[1].permission || {};
      renderMenuSelect(form.elements.menu_id, items[0].tree || [], perm.menu_id || 0);
      fillForm(form, perm);
      form.elements.menu_id.setAttribute("data-checked-id", perm.parent_id || 0);
      return loadPermissionParents(form, perm.parent_id || 0);
    }).catch(function (err) { setAlert(err.message, "danger"); });
    form.addEventListener("submit", function (e) {
      e.preventDefault();
      var parentIDs = checkedValues('#permission-parent-tree input[name="parent_id"]');
      var payload = {
        parent_id: parentIDs[0] || 0,
        menu_id: formNumber(form, "menu_id"),
        name: form.elements.name.value.trim(),
        code: form.elements.code.value.trim(),
        method: form.elements.method.value,
        path: form.elements.path.value.trim(),
        sort: formNumber(form, "sort"),
        status: formBool(form, "status")
      };
      api(mode === "edit" ? "/api/permissions/" + id : "/api/permissions", {
        method: mode === "edit" ? "PUT" : "POST",
        body: JSON.stringify(payload)
      }).then(function () { window.location.href = "/permissions?flash=" + encodeURIComponent("保存成功"); })
        .catch(function (err) { setAlert(err.message, "danger"); });
    });
  }

  function renderMenuSelect(select, nodes, selected) {
    select.innerHTML = '<option value="0">不绑定菜单</option>';
    flattenTree(nodes, function (node) { return node.menu || {}; }).forEach(function (row) {
      var menu = row.item;
      var id = Number(value(menu, "id", "ID"));
      var option = document.createElement("option");
      option.value = id;
      option.innerHTML = indent(row.level) + text(value(menu, "title", "Title"));
      option.selected = Number(selected) === id;
      select.appendChild(option);
    });
  }

  function loadPermissionParents(form, checkedID) {
    var select = form.elements.menu_id;
    var params = new URLSearchParams();
    params.set("menu_id", select.value || "0");
    params.set("current_id", select.getAttribute("data-current-id") || "0");
    params.set("checked", checkedID || 0);
    return api("/api/permission-parents?" + params.toString()).then(function (payload) {
      document.querySelector("#permission-parent-tree").innerHTML = renderPermissionParentChecks(payload.tree || [], Number(checkedID || 0));
    });
  }

  function appendIDs(params, key, ids) {
    if (!ids || !ids.length) {
      params.append(key, "");
      return;
    }
    ids.forEach(function (id) { params.append(key, id); });
  }

  function collectPermissionIDs() {
    var ids = checkedValues('#role-permission-tree input[name="permission_ids"]');
    document.querySelectorAll('#role-hidden-permissions input[name="permission_ids"]').forEach(function (input) {
      ids.push(Number(input.value));
    });
    return Array.from(new Set(ids.filter(Boolean)));
  }

  function syncHiddenSelected(visibleSelector, hiddenSelector, name, selectedIDs) {
    var visible = new Set(Array.prototype.slice.call(document.querySelectorAll(visibleSelector)).map(function (input) {
      return Number(input.value);
    }));
    var box = document.querySelector(hiddenSelector);
    if (!box) return;
    box.innerHTML = "";
    (selectedIDs || []).forEach(function (id) {
      id = Number(id);
      if (!id || visible.has(id)) return;
      var input = document.createElement("input");
      input.type = "hidden";
      input.name = name;
      input.value = id;
      box.appendChild(input);
    });
  }

  function renderMenuChecks(nodes) {
    if (!nodes.length) return '<div class="empty-cell">暂无菜单</div>';
    return '<ul class="check-tree">' + nodes.map(function (node) {
      var menu = node.menu || {};
      var children = node.children || [];
      var id = value(menu, "id", "ID");
      var path = value(menu, "path", "Path");
      var toggle = children.length ? '<button class="tree-toggle" type="button" data-tree-toggle>▾</button>' : '<span class="tree-spacer"></span>';
      return '<li><div class="tree-line">' + toggle +
        '<label><input type="checkbox" name="menu_ids" value="' + id + '"' + (node.checked ? " checked" : "") + ">" +
        "<span>" + text(value(menu, "title", "Title")) + "</span>" + (path ? "<code>" + text(path) + "</code>" : "") + "</label></div>" +
        (children.length ? '<div class="tree-children">' + renderMenuChecks(children) + "</div>" : "") + "</li>";
    }).join("") + "</ul>";
  }

  function renderPermissionChecks(nodes) {
    if (!nodes.length) return '<div class="empty-cell">暂无权限</div>';
    return '<ul class="check-tree">' + nodes.map(function (node) {
      var perm = node.permission || {};
      var children = node.children || [];
      var id = value(perm, "id", "ID");
      var toggle = children.length ? '<button class="tree-toggle" type="button" data-tree-toggle>▾</button>' : '<span class="tree-spacer"></span>';
      return '<li><div class="tree-line">' + toggle +
        '<label><input type="checkbox" name="permission_ids" value="' + id + '"' + (node.checked ? " checked" : "") + ">" +
        "<span>" + text(value(perm, "name", "Name")) + "</span><code>" + text(value(perm, "code", "Code")) + "</code></label></div>" +
        (children.length ? '<div class="tree-children">' + renderPermissionChecks(children) + "</div>" : "") + "</li>";
    }).join("") + "</ul>";
  }

  function renderPermissionParentChecks(nodes, checkedID) {
    var topChecked = Number(checkedID || 0) === 0 ? "" : "";
    var top = '<ul class="check-tree"><li><div class="tree-line"><span class="tree-spacer"></span><label>' +
      '<input class="single-parent-check" type="checkbox" name="parent_id" value="0"' + topChecked + '><span>作为顶级权限</span>' +
      '</label></div></li>';
    return top + (nodes || []).map(function render(node) {
      var perm = node.permission || {};
      var children = node.children || [];
      var id = value(perm, "id", "ID");
      var toggle = children.length ? '<button class="tree-toggle" type="button" data-tree-toggle>▾</button>' : '<span class="tree-spacer"></span>';
      return '<li><div class="tree-line">' + toggle +
        '<label><input class="single-parent-check" type="checkbox" name="parent_id" value="' + id + '"' + (node.checked ? " checked" : "") + ">" +
        "<span>" + text(value(perm, "name", "Name")) + "</span><code>" + text(value(perm, "code", "Code")) + "</code></label></div>" +
        (children.length ? '<div class="tree-children"><ul class="check-tree">' + children.map(render).join("") + "</ul></div>" : "") + "</li>";
    }).join("") + "</ul>";
  }

  on("[data-sidebar-toggle]", "click", function () {
    document.body.classList.toggle("sidebar-collapsed");
  });

  on("[data-logout]", "click", function () {
    api("/api/logout", { method: "POST" }).finally(function () {
      window.location.href = "/login";
    });
  });

  on("[data-tree-toggle]", "click", function (e, target) {
    var branch = target.closest(".nav-branch");
    if (branch) {
      branch.classList.toggle("open");
      return;
    }
    var line = target.closest(".tree-line");
    if (!line) return;
    var children = line.parentElement.querySelector(":scope > .tree-children");
    if (children) {
      children.classList.toggle("collapsed");
      target.textContent = children.classList.contains("collapsed") ? "▸" : "▾";
    }
  });

  on(".single-parent-check", "change", function (e, target) {
    if (!target.checked) return;
    document.querySelectorAll(".single-parent-check").forEach(function (input) {
      if (input !== target) input.checked = false;
    });
  });

  on(".tree-table-caret", "click", function (e, target) {
    var row = target.closest("tr");
    if (!row) return;
    var level = Number(row.getAttribute("data-level") || "0");
    var collapse = !row.classList.contains("collapsed");
    row.classList.toggle("collapsed", collapse);
    var next = row.nextElementSibling;
    while (next && Number(next.getAttribute("data-level") || "0") > level) {
      if (collapse) {
        next.hidden = true;
      } else if (Number(next.getAttribute("data-level") || "0") === level + 1) {
        next.hidden = false;
      }
      next = next.nextElementSibling;
    }
  });

  on("[data-expand-table]", "click", function () {
    document.querySelectorAll("[data-tree-table] tr").forEach(function (row) {
      row.hidden = false;
      row.classList.remove("collapsed");
    });
  });

  on("[data-collapse-table]", "click", function () {
    document.querySelectorAll("[data-tree-table] tbody tr").forEach(function (row) {
      var level = Number(row.getAttribute("data-level") || "0");
      row.hidden = level > 0;
      row.classList.toggle("collapsed", level === 0);
    });
  });

  on("[data-check-all]", "click", function (e, target) {
    document.querySelectorAll(target.getAttribute("data-check-all") + " input[type=checkbox]").forEach(function (input) {
      input.checked = true;
    });
  });

  on("[data-uncheck-all]", "click", function (e, target) {
    var selector = target.getAttribute("data-uncheck-all");
    document.querySelectorAll(selector + " input[type=checkbox]").forEach(function (input) {
      input.checked = false;
    });
    if (selector === "#role-permission-tree") {
      var hidden = document.querySelector("#role-hidden-permissions");
      if (hidden) hidden.innerHTML = "";
    }
  });

  on("[data-role-access-filter]", "change", function (e, select) {
    loadRoleAccess(select.getAttribute("data-role-id") || "0", checkedValues('#role-menu-tree input[name="menu_ids"]'), collectPermissionIDs())
      .catch(function (err) { setAlert(err.message, "danger"); });
  });

  on("[data-permission-menu-select]", "change", function (e, select) {
    var form = select.closest("form");
    var checked = checkedValues('#permission-parent-tree input[name="parent_id"]')[0] || 0;
    loadPermissionParents(form, checked).catch(function (err) { setAlert(err.message, "danger"); });
  });

  on("[data-delete-user]", "click", function (e, target) {
    if (!confirm("确认删除该用户？")) return;
    api("/api/users/" + target.getAttribute("data-delete-user"), { method: "DELETE" }).then(loadUsers).catch(function (err) { setAlert(err.message, "danger"); });
  });

  on("[data-delete-role]", "click", function (e, target) {
    if (!confirm("确认删除该角色？")) return;
    api("/api/roles/" + target.getAttribute("data-delete-role"), { method: "DELETE" }).then(loadRoles).catch(function (err) { setAlert(err.message, "danger"); });
  });

  on("[data-delete-menu]", "click", function (e, target) {
    if (!confirm("确认删除该菜单？子菜单会移动到顶级。")) return;
    api("/api/menus/" + target.getAttribute("data-delete-menu"), { method: "DELETE" }).then(loadMenus).catch(function (err) { setAlert(err.message, "danger"); });
  });

  on("[data-delete-permission]", "click", function (e, target) {
    if (!confirm("确认删除该权限？子权限会移动到顶级。")) return;
    api("/api/permissions/" + target.getAttribute("data-delete-permission"), { method: "DELETE" }).then(loadPermissions).catch(function (err) { setAlert(err.message, "danger"); });
  });

  initLogin();
  initShell();
  initDashboard();
  initUsersIndex();
  initUserForm();
  initRolesIndex();
  initRoleForm();
  initMenusIndex();
  initMenuForm();
  initPermissionsIndex();
  initPermissionForm();
})();
