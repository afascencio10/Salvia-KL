$(document).ready(function(){
    $(".dropdown-menu a").click(function(){
        var selectedValue = $(this).data("value");
        $("#estado-dropdown-btn").text($(this).text());
        $("#estado-dropdown-btn").val(selectedValue);
    });
});
