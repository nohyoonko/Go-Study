(function ($) {
    'use strict';
    $(function () {
        var todoListItem = $('.todo-list');
        var todoListInput = $('.todo-list-input');
        //add를 클릭하면 불리는 function으로 server에 add 요청
        $('.todo-list-add-btn').on("click", function (event) {
            event.preventDefault();

            var item = $(this).prevAll('.todo-list-input').val();

            if (item) {
                //server에 add 요청, server가 json object로 줌
                $.post("/todos", {name:item}, addItem)
                //todoListItem.append("<li><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
                todoListInput.val("");
            }
        });

        var addItem = function(item) {
            if (item.completed) {
                todoListItem.append("<li class='completed'"+" id='"+item.id+"'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' checked='checked'/>" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            }
            else{
                todoListItem.append("<li "+" id='"+item.id+"'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            }
        };

        //todos 경로를 호출하면 데이터가 옴
        $.get('/todos', function(items) {
            items.forEach(e => {
                addItem(e)
            });
        });

        todoListItem.on('change', '.checkbox', function () {
            var id = $(this).closest("li").attr('id');
            var $self = $(this);
            var complete = true;

            if ($self.attr('checked')) {
                complete = false;
            }
            // 응답(data)이 오면 실행, 토글되니까 알려줘야한다
            $.get("complete-todo/"+id+"?complete="+complete, function(data){
                if (complete) {
                    $self.attr('checked', 'checked');
                } else {
                    $self.removeAttr('checked');
                }

                $self.closest("li").toggleClass('completed');
            })
        });

        todoListItem.on('click', '.remove', function () {
            // url: todos/id    method: DELETE
            //this = .remove 에 가장 가까운 li를 가져오면 id를 가져올 수 있음
            var id = $(this).closest("li").attr('id');
            var $self = $(this);
            $.ajax({
                url: "todos/" + id,
                type: "DELETE",
                success: function(data) {
                    if (data.success) {
                        $self.parent().remove();
                    }
                }
            })
            //$(this).parent().remove();
        });

    });
})(jQuery);