 /**
	*��ǩ�ؼ�
    *���ܣ���Enter��Tab��ʧȥ����ȷ����ǩ������ϣ�˫�����ֿ��Ա༭�ñ�ǩ��������棨������ʾɾ���ñ�ǩ
    *tabControl:function
    *����˵����
    *initTabCount:int һ��ʼ��ʼ����ǩ������������
    *maxTabCount:int �����ɽ������ı�ǩ������
    *tabMaxLen:int ÿ����ǩ������������ַ����ȣ�
    *tabW:int ��ǩ�����Ŀ�ȣ�
    *tabH:int ��ǩ�����ĸ߶ȣ�
    *tipTOffset:int ��ʾ��Ϣ���ǩ������topƫ������
    *tipLOffset:int ��ʾ��Ϣ���ǩ������leftƫ������
    *tags:string ��ʼ���ı�ǩ���ݣ��Զ���Ϊ�����
 **/
$.fn.extend({
	tabControl: function(options, tags) {
		var defOpt = {
			initTabCount: 1,
			maxTabCount: 10,
			tabMaxLen: 10,
			tabW: 150,
			tabH: 15,
			tipTOffset: 5,
			tipLOffset: 0
		};
		var opts = $.extend(defOpt, options);
		var _tags = [];
		if (tags) {
			tags = tags.replace(/[^A-Za-z0-9_,\u4E00-\u9FA5]+/gi, "").replace(/^,+|,+$/gi, "");//������Ӣ�ġ����֡��»�˿�����ŵ������ַ���ȥ�����Ҳ����Զ��ſ�ͷ�����
			_tags = tags.split(',');
		}
		_tags = _tags.length > opts.maxTabCount ? _tags.slice(0, opts.maxTabCount - 1) : _tags;
		opts.initTabCount = opts.maxTabCount <= _tags.length ? _tags.length: _tags.length + (opts.maxTabCount - _tags.length > opts.initTabCount ? opts.initTabCount: opts.maxTabCount - _tags.length);
		var checkReg = /[^A-Za-z0-9_\u4E00-\u9FA5]+/gi;//ƥ��Ƿ��ַ�
		var initTab = function(obj, index) {//��ʼ����ǩ����
			var textHtml = "<input class='tabinput' name='tabinput' style='width:" + opts.tabW + "px;height:" + opts.tabH + "px;' type='text'/>";
			obj.append(textHtml);
			if (_tags[index]) {
				var __inputobj = $("input[type='text'][name='tabinput']", obj).eq(index);
				__inputobj.val(_tags[index].substr(0, opts.tabMaxLen)).css("display", "none");
				compTab(obj, __inputobj, _tags[index].substr(0, opts.tabMaxLen));
			}
			$("input[type='text'][name='tabinput']:last", obj).bind("keydown blur click",
			function(event) {
				if (event.type == "click") {//����¼�����������յ�һ���¼�����(event)������ͨ��������ֹ���������Ĭ�ϵ���Ϊ���������ȡ��Ĭ�ϵ���Ϊ��event.preventDefault()����������ֹ�¼����ݡ�event.stopPropagation()��������¼����������뷵��false��
					return false;
				}
				if (event.keyCode == 13 || event.keyCode == 9 || event.type == "blur") {
					event.preventDefault();//��Ҫ�Ǟ���tab�I����Ҫ׌��ǰԪ��ʧȥ���c������ֹ���������Ĭ�ϵ���Ϊ��
					event.stopPropagation();
					var inputObj = $(this);
					var value = $(this).val().replace(/\s+/gi, "");
					if ((event.keyCode == 13 || event.keyCode == 9) && value != "")//��Ҫ��̎��IE
					 inputObj.data("isIEKeyDown", true);
					if (event.type == "blur" && inputObj.data("isIEKeyDown")) {
						inputObj.removeData("isIEKeyDown");
						return;
					}
					if (value != "") {
						if (value.length > opts.tabMaxLen) {
							showMes($(this), "������1��" + opts.tabMaxLen + "���ַ����ȵı�ǩ");
							return;
						}
						var _match = value.match(checkReg);
						if (!_match) {
							compTab(obj, inputObj, value);
							if ($("input[type='text'][name='tabinput']", obj).length < opts.maxTabCount) {
								if (!inputObj.data("isModify"))
								 initTab(obj);
								else if (!$("input[type='text'][name='tabinput']", obj).is(":hidden")) {
									initTab(obj);
								}
							}
							$("input[type='text']:last", obj).focus();
							hideErr();
						}
						 else {
							showMes(inputObj, "���ݲ��ܰ����Ƿ��ַ���{0}����".replace("{0}", _match.join(" ")));
						}
					}
					 else {
						if (event.type != "blur")
						 showMes(inputObj, "���ݲ����");
					}
				}
			}).bind("focus",
			function() {
				hideErr();
			});
		}
		//��ɱ�ǩ��д
		var compTab = function(obj, inputObj, value) {
			inputObj.next("span").remove();//ɾ������inputԪ�����span
			var _span = "<span name='tab' id='radius'><b>" + value + "</b><a id='deltab'>��</a></span>";
			inputObj.after(_span).hide();
			inputObj.next("span").find("a").click(function() {
				if (confirm("ȷ��ɾ���ñ�ǩ��")) {
					inputObj.next("span").remove();
					inputObj.remove();
					if ($("span[name='tab']", obj).length == opts.maxTabCount - 1)
					 initTab(obj);
				}
			});
			inputObj.next("span").dblclick(function() {
				inputObj.data("isModify", true).next("span").remove();
				inputObj.show().focus();
			});
		}
		return this.each(function() {
			var jqObj = $(this);
			for (var i = 0; i < opts.initTabCount; i++) {
				initTab(jqObj, i);
			}
			jqObj.data("isInit", true);
			jqObj.click(function() {
				$("input[type='text'][name='tabinput']", jqObj).each(function() {
					if ($(this).val() == "") {
						$(this).focus();
						return false;
					}
				});
			});
		});
		function showMes(inputObj, mes) {
			var _offset = inputObj.offset();
			var _mesHtml = "<div id='errormes' class='radius_shadow' style='position:absolute;left:" + (_offset.left + opts.tipLOffset) + "px;top:" + (_offset.top + opts.tabH + opts.tipTOffset) + "px;'>" + mes + "</div>";
			$("#errormes").remove();
			$("body").append(_mesHtml);
		}
		function hideErr() {
			$("#errormes").hide();
		}
		function showErr() {
			$("#errormes").show();
		}
	},
	getTabVals: function() {//��ȡ��ǰ���������ɵ�tabֵ�������һά����
		var obj = $(this);
		var values = [];
		obj.children("span[name=\"tab\"][id^=\"radius\"]").find("b").text(function(index, text) {
			var checkReg = /[^A-Za-z0-9_\u4E00-\u9FA5]+/gi;
			values.push(text.replace(checkReg, ""));
		});
		return values;
	}
});