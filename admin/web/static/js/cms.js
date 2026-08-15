(function () {
  var root = document.querySelector("[data-cms-page]");
  if (!root) return;

  var resource = root.getAttribute("data-resource");
  var meta = null;
  var page = 1;
  var pageSize = 20;
  var query = "";
  var editingID = 0;
  var relationCache = {};
  var dialog = document.querySelector("[data-cms-dialog]");
  var form = document.querySelector("[data-cms-form]");

  function api(path, options) {
    options = options || {};
    var headers = options.headers || {};
    if (options.body && !(options.body instanceof FormData)) headers["Content-Type"] = "application/json";
    return fetch(path, Object.assign({}, options, { headers: headers })).then(function (response) {
      return response.json().catch(function () { return {}; }).then(function (payload) {
        if (response.status === 401) window.location.href = "/login";
        if (!response.ok) throw new Error(payload.error || "请求失败");
        return payload;
      });
    });
  }

  function escapeHTML(value) {
    return String(value == null ? "" : value)
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;");
  }

  function notify(message, danger) {
    var old = document.querySelector(".runtime-alert");
    if (old) old.remove();
    var el = document.createElement("div");
    el.className = "alert runtime-alert " + (danger ? "danger" : "success");
    el.textContent = message;
    document.querySelector(".content").prepend(el);
    setTimeout(function () { el.remove(); }, 3200);
  }

  function loadMeta() {
    return api("/api/cms/" + resource + "/meta").then(function (payload) {
      meta = payload.resource;
      meta.abilities = payload.abilities || {};
      if (!meta.abilities.create) document.querySelector("[data-cms-create]").hidden = true;
      renderHead();
      return loadRelations();
    });
  }

  function loadRelations() {
    var resources = Array.from(new Set((meta.form_fields || []).map(function (field) {
      return field.relation;
    }).filter(Boolean)));
    return Promise.all(resources.map(function (name) {
      return api("/api/cms/" + name + "?page_size=500").then(function (payload) {
        relationCache[name] = payload.items || [];
      }).catch(function () {
        relationCache[name] = [];
      });
    }));
  }

  function renderHead() {
    document.querySelector("[data-cms-head]").innerHTML = "<tr><th>ID</th>" +
      (meta.list_fields || []).map(function (field) {
        return "<th>" + escapeHTML(field.label) + "</th>";
      }).join("") + '<th class="actions">操作</th></tr>';
  }

  function loadRows() {
    var params = new URLSearchParams({ page: page, page_size: pageSize, q: query });
    document.querySelector("[data-cms-body]").innerHTML = '<tr><td class="empty-cell" colspan="' +
      ((meta.list_fields || []).length + 2) + '">加载中...</td></tr>';
    api("/api/cms/" + resource + "?" + params).then(function (payload) {
      renderRows(payload.items || []);
      renderPagination(payload.total || 0, payload.page || 1, payload.page_size || pageSize);
    }).catch(function (error) {
      notify(error.message, true);
    });
  }

  function renderRows(items) {
    var body = document.querySelector("[data-cms-body]");
    if (!items.length) {
      body.innerHTML = '<tr><td class="empty-cell" colspan="' +
        ((meta.list_fields || []).length + 2) + '">暂无数据</td></tr>';
      return;
    }
    body.innerHTML = items.map(function (item) {
      var cells = (meta.list_fields || []).map(function (field) {
        return "<td>" + renderValue(item[field.name], field.type) + "</td>";
      }).join("");
      var edit = meta.abilities.update ? '<button class="btn small" type="button" data-cms-edit="' + item.id + '">编辑</button>' : "";
      var remove = meta.abilities.delete ? '<button class="btn small danger" type="button" data-cms-delete="' + item.id + '">删除</button>' : "";
      return "<tr><td>" + item.id + "</td>" + cells + '<td class="actions">' + edit + remove +
        (edit || remove ? "" : "-") + "</td></tr>";
    }).join("");
  }

  function renderValue(value, type) {
    if (type === "boolean") return '<span class="badge ' + (value ? "success" : "secondary") + '">' + (value ? "启用" : "禁用") + "</span>";
    if (type === "image") return value ? '<img class="cms-thumb" src="' + escapeHTML(value) + '" alt="">' : "-";
    if (type === "code") return "<code>" + escapeHTML(value || "-") + "</code>";
    if (type === "badge") return '<span class="badge primary">' + escapeHTML(value || "-") + "</span>";
    if (type === "datetime") return value ? escapeHTML(new Date(value).toLocaleString("zh-CN", { hour12: false })) : "-";
    if (type === "truncate") return '<span class="cms-truncate" title="' + escapeHTML(value || "") + '">' + escapeHTML(value || "-") + "</span>";
    return escapeHTML(value == null || value === "" ? "-" : value);
  }

  function renderPagination(total, current, size) {
    var pages = Math.max(1, Math.ceil(total / size));
    var box = document.querySelector("[data-cms-pagination]");
    if (pages <= 1) {
      box.innerHTML = '<span>共 ' + total + " 条</span>";
      return;
    }
    var start = Math.max(1, Math.min(pages - 6, current - 3));
    var buttons = [];
    for (var i = start; i <= Math.min(pages, start + 6); i++) {
      buttons.push('<button class="btn small ' + (i === current ? "primary" : "") +
        '" data-cms-page-number="' + i + '">' + i + "</button>");
    }
    box.innerHTML = '<button class="btn small" data-cms-page-number="1">首页</button>' + buttons.join("") +
      '<button class="btn small" data-cms-page-number="' + pages + '">末页</button><span>共 ' + total + " 条</span>";
  }

  function openForm(item) {
    item = item || {};
    editingID = Number(item.id || 0);
    document.querySelector("[data-cms-dialog-title]").textContent = editingID ? "编辑" + meta.title : "新增" + meta.title;
    document.querySelector("[data-cms-fields]").innerHTML = (meta.form_fields || []).map(function (field) {
      return renderField(field, item[field.name]);
    }).join("");
    syncSettingImageTools();
    if (typeof dialog.showModal === "function") dialog.showModal();
    else dialog.setAttribute("open", "open");
  }

  function renderField(field, value) {
    var classes = "field" + (field.wide ? " wide" : "");
    var required = field.required ? " required" : "";
    var disabled = field.read_only ? " disabled" : "";
    var help = field.help ? "<small>" + escapeHTML(field.help) + "</small>" : "";
    var control = "";

    if (field.type === "boolean") {
      var checked = value == null || value ? " checked" : "";
      return '<label class="switch-row cms-switch"><input type="checkbox" name="' + field.name + '"' + checked + disabled +
        "><span>" + escapeHTML(field.label) + "</span></label>";
    }
    if (field.type === "images") return renderGalleryField(field, value || []);
    if (field.type === "image" || field.type === "image_value") return renderImageField(field, value || "", required, disabled, help, classes);

    if (field.type === "textarea" || field.type === "richtext") {
      control = '<textarea name="' + field.name + '" rows="' + (field.type === "richtext" ? "13" : "5") + '"' +
        required + disabled + ">" + escapeHTML(value || "") + "</textarea>";
    } else if (field.type === "select") {
      control = '<select name="' + field.name + '"' + required + disabled + ">" + (field.options || []).map(function (option) {
        return '<option value="' + escapeHTML(option.value) + '"' +
          (String(value || "") === String(option.value) ? " selected" : "") + ">" + escapeHTML(option.label) + "</option>";
      }).join("") + "</select>";
    } else if (field.type === "relation") {
      var options = relationCache[field.relation] || [];
      control = '<select name="' + field.name + '"' + required + disabled + '><option value="0">无 / 顶级</option>' +
        options.map(function (option) {
          var id = option[field.relation_id || "id"];
          return '<option value="' + id + '"' + (Number(value || 0) === Number(id) ? " selected" : "") + ">" +
            escapeHTML(option[field.relation_text || "name"]) + " (#" + id + ")</option>";
        }).join("") + "</select>";
    } else {
      var type = field.type === "number" ? "number" : field.type === "datetime" ? "datetime-local" : "text";
      var formatted = field.type === "datetime" && value ? new Date(value).toISOString().slice(0, 16) : (value == null ? "" : value);
      control = '<input type="' + type + '" name="' + field.name + '" value="' + escapeHTML(formatted) + '"' + required + disabled + ">";
    }
    return '<label class="' + classes + '"><span>' + escapeHTML(field.label) + "</span>" + control + help + "</label>";
  }

  function renderImageField(field, value, required, disabled, help, classes) {
    var input = field.type === "image_value"
      ? '<textarea name="' + field.name + '" rows="5" data-image-url' + required + disabled + ">" + escapeHTML(value) + "</textarea>"
      : '<input type="text" name="' + field.name + '" value="' + escapeHTML(value) + '" data-image-url placeholder="图片 URL"' + required + disabled + ">";
    var tools = field.read_only ? "" : '<div class="cms-image-tools" data-image-tools>' +
      '<input type="file" accept="image/jpeg,image/png,image/gif,image/webp" data-image-file>' +
      '<button class="btn small primary" type="button" data-image-upload>上传图片</button>' +
      '<button class="btn small" type="button" data-image-clear>清除</button></div>';
    return '<label class="' + classes + '" data-image-control data-image-kind="' + field.type + '"><span>' +
      escapeHTML(field.label) + "</span>" + input + tools + '<div class="cms-image-preview" data-image-preview>' +
      (value ? '<img src="' + escapeHTML(value) + '" alt="">' : "") + "</div>" + help + "</label>";
  }

  function renderGalleryField(field, values) {
    return '<div class="field wide cms-gallery" data-gallery-field="' + field.name + '"><span>' + escapeHTML(field.label) + "</span>" +
      '<div class="cms-gallery-list" data-gallery-list>' + values.map(renderGalleryItem).join("") + "</div>" +
      '<div class="cms-gallery-upload"><input type="file" accept="image/jpeg,image/png,image/gif,image/webp" multiple data-gallery-files>' +
      '<button class="btn small primary" type="button" data-gallery-upload>批量上传</button></div>' +
      '<div class="cms-gallery-url"><input type="text" data-gallery-new-url placeholder="也可填写已有图片 URL">' +
      '<button class="btn small" type="button" data-gallery-add-url>添加 URL</button></div>' +
      (field.help ? "<small>" + escapeHTML(field.help) + "</small>" : "") + "</div>";
  }

  function renderGalleryItem(item) {
    var url = typeof item === "string" ? item : (item.image_url || "");
    var alt = typeof item === "string" ? "" : (item.alt_text || "");
    return '<div class="cms-gallery-item" data-gallery-item><div class="cms-gallery-thumb" data-gallery-preview>' +
      (url ? '<img src="' + escapeHTML(url) + '" alt="">' : "") + "</div>" +
      '<div class="cms-gallery-fields"><input type="text" value="' + escapeHTML(url) + '" data-gallery-url placeholder="图片 URL">' +
      '<input type="text" value="' + escapeHTML(alt) + '" data-gallery-alt placeholder="替代文本"></div>' +
      '<div class="cms-gallery-actions"><button class="btn small" type="button" data-gallery-up title="上移">↑</button>' +
      '<button class="btn small" type="button" data-gallery-down title="下移">↓</button>' +
      '<button class="btn small danger" type="button" data-gallery-remove>移除</button></div></div>';
  }

  function collectForm() {
    var data = {};
    (meta.form_fields || []).forEach(function (field) {
      if (field.read_only) return;
      if (field.type === "images") {
        var gallery = form.querySelector('[data-gallery-field="' + field.name + '"]');
        data[field.name] = Array.prototype.slice.call(gallery.querySelectorAll("[data-gallery-item]")).map(function (item) {
          return {
            image_url: item.querySelector("[data-gallery-url]").value.trim(),
            alt_text: item.querySelector("[data-gallery-alt]").value.trim()
          };
        }).filter(function (item) { return item.image_url; });
        return;
      }
      var input = form.elements[field.name];
      if (!input) return;
      if (field.type === "boolean") data[field.name] = input.checked;
      else if (field.type === "number" || field.type === "relation") data[field.name] = Number(input.value || 0);
      else data[field.name] = input.value;
    });
    return data;
  }

  function uploadImage(file) {
    var body = new FormData();
    body.append("file", file);
    return api("/api/cms/" + resource + "/upload", { method: "POST", body: body }).then(function (payload) {
      return payload.item.url;
    });
  }

  function updateImagePreview(control) {
    var preview = control.querySelector("[data-image-preview]");
    var url = control.querySelector("[data-image-url]").value.trim();
    preview.innerHTML = "";
    if (!url) return;
    var image = document.createElement("img");
    image.src = url;
    image.alt = "";
    preview.appendChild(image);
  }

  function previewSelectedImage(control, file) {
    var preview = control.querySelector("[data-image-preview]");
    var objectURL = URL.createObjectURL(file);
    preview.hidden = false;
    preview.innerHTML = "";
    var image = document.createElement("img");
    image.src = objectURL;
    image.alt = file.name;
    image.onload = function () { URL.revokeObjectURL(objectURL); };
    preview.appendChild(image);
  }

  function updateGalleryPreview(item) {
    var preview = item.querySelector("[data-gallery-preview]");
    var url = item.querySelector("[data-gallery-url]").value.trim();
    preview.innerHTML = "";
    if (!url) return;
    var image = document.createElement("img");
    image.src = url;
    image.alt = "";
    preview.appendChild(image);
  }

  function syncSettingImageTools() {
    var valueType = form.elements.value_type;
    form.querySelectorAll('[data-image-kind="image_value"]').forEach(function (control) {
      var hidden = !!valueType && valueType.value !== "image";
      var tools = control.querySelector("[data-image-tools]");
      var preview = control.querySelector("[data-image-preview]");
      if (tools) tools.hidden = hidden;
      if (preview) preview.hidden = hidden;
    });
  }

  document.querySelector("[data-cms-create]").addEventListener("click", function () {
    openForm({ status: true, sort: 10, published_at: new Date().toISOString() });
  });

  document.querySelectorAll("[data-cms-close]").forEach(function (button) {
    button.addEventListener("click", function () { dialog.close(); });
  });

  document.querySelector("[data-cms-search]").addEventListener("submit", function (event) {
    event.preventDefault();
    query = event.currentTarget.elements.q.value.trim();
    page = 1;
    loadRows();
  });

  root.addEventListener("click", function (event) {
    var edit = event.target.closest("[data-cms-edit]");
    if (edit) {
      api("/api/cms/" + resource + "/" + edit.getAttribute("data-cms-edit")).then(function (payload) {
        openForm(payload.item);
      }).catch(function (error) { notify(error.message, true); });
      return;
    }
    var remove = event.target.closest("[data-cms-delete]");
    if (remove) {
      if (!confirm("确认删除这条记录？")) return;
      api("/api/cms/" + resource + "/" + remove.getAttribute("data-cms-delete"), { method: "DELETE" }).then(function (payload) {
        notify(payload.message);
        loadRows();
      }).catch(function (error) { notify(error.message, true); });
      return;
    }
    var pager = event.target.closest("[data-cms-page-number]");
    if (pager) {
      page = Number(pager.getAttribute("data-cms-page-number"));
      loadRows();
    }
  });

  form.addEventListener("click", function (event) {
    var upload = event.target.closest("[data-image-upload]");
    if (upload) {
      var control = upload.closest("[data-image-control]");
      var file = control.querySelector("[data-image-file]").files[0];
      if (!file) {
        notify("请先选择图片文件", true);
        return;
      }
      upload.disabled = true;
      uploadImage(file).then(function (url) {
        control.querySelector("[data-image-url]").value = url;
        updateImagePreview(control);
        notify("图片上传成功");
      }).catch(function (error) {
        notify(error.message, true);
      }).finally(function () {
        upload.disabled = false;
      });
      return;
    }

    var clear = event.target.closest("[data-image-clear]");
    if (clear) {
      var imageControl = clear.closest("[data-image-control]");
      imageControl.querySelector("[data-image-url]").value = "";
      imageControl.querySelector("[data-image-file]").value = "";
      updateImagePreview(imageControl);
      return;
    }

    var galleryUpload = event.target.closest("[data-gallery-upload]");
    if (galleryUpload) {
      var gallery = galleryUpload.closest("[data-gallery-field]");
      var files = Array.prototype.slice.call(gallery.querySelector("[data-gallery-files]").files);
      if (!files.length) {
        notify("请先选择图片文件", true);
        return;
      }
      if (gallery.querySelectorAll("[data-gallery-item]").length + files.length > 50) {
        notify("产品图集最多保留 50 张图片", true);
        return;
      }
      galleryUpload.disabled = true;
      Promise.all(files.map(uploadImage)).then(function (urls) {
        gallery.querySelector("[data-gallery-list]").insertAdjacentHTML("beforeend", urls.map(function (url) {
          return renderGalleryItem({ image_url: url, alt_text: "" });
        }).join(""));
        gallery.querySelector("[data-gallery-files]").value = "";
        notify("已上传 " + urls.length + " 张图片");
      }).catch(function (error) {
        notify(error.message, true);
      }).finally(function () {
        galleryUpload.disabled = false;
      });
      return;
    }

    var addURL = event.target.closest("[data-gallery-add-url]");
    if (addURL) {
      var urlField = addURL.closest("[data-gallery-field]").querySelector("[data-gallery-new-url]");
      var url = urlField.value.trim();
      if (!url) return;
      addURL.closest("[data-gallery-field]").querySelector("[data-gallery-list]")
        .insertAdjacentHTML("beforeend", renderGalleryItem({ image_url: url, alt_text: "" }));
      urlField.value = "";
      return;
    }

    var galleryItem = event.target.closest("[data-gallery-item]");
    if (!galleryItem) return;
    if (event.target.closest("[data-gallery-remove]")) galleryItem.remove();
    else if (event.target.closest("[data-gallery-up]") && galleryItem.previousElementSibling) {
      galleryItem.parentNode.insertBefore(galleryItem, galleryItem.previousElementSibling);
    } else if (event.target.closest("[data-gallery-down]") && galleryItem.nextElementSibling) {
      galleryItem.parentNode.insertBefore(galleryItem.nextElementSibling, galleryItem);
    }
  });

  form.addEventListener("input", function (event) {
    var imageInput = event.target.closest("[data-image-url]");
    if (imageInput) updateImagePreview(imageInput.closest("[data-image-control]"));
    var galleryInput = event.target.closest("[data-gallery-url]");
    if (galleryInput) updateGalleryPreview(galleryInput.closest("[data-gallery-item]"));
  });

  form.addEventListener("change", function (event) {
    var imageFile = event.target.closest("[data-image-file]");
    if (imageFile && imageFile.files[0]) {
      previewSelectedImage(imageFile.closest("[data-image-control]"), imageFile.files[0]);
    }
    if (event.target.name === "value_type") syncSettingImageTools();
  });

  form.addEventListener("submit", function (event) {
    event.preventDefault();
    var url = "/api/cms/" + resource + (editingID ? "/" + editingID : "");
    var submit = form.querySelector('button[type="submit"]');
    submit.disabled = true;
    api(url, {
      method: editingID ? "PUT" : "POST",
      body: JSON.stringify(collectForm())
    }).then(function (payload) {
      dialog.close();
      notify(payload.message || "保存成功");
      loadRows();
    }).catch(function (error) {
      notify(error.message, true);
    }).finally(function () {
      submit.disabled = false;
    });
  });

  loadMeta().then(loadRows).catch(function (error) { notify(error.message, true); });
})();
