(function () {
  'use strict';
  window.addEventListener('load', function () {
    // Fetch all the forms we want to apply custom Bootstrap validation styles to
    var forms = document.getElementsByClassName('needs-validation');
    // Loop over them and prevent submission
    var validation = Array.prototype.filter.call(forms, function (form) {
      form.addEventListener('submit', function (event) {
        if (form.checkValidity() === false) {
          event.preventDefault();
          event.stopPropagation();
        }
        form.classList.add('was-validated');
      }, false);
    });
  }, false);
})();
$('document').ready(function () {
  $('#add-pic').click(function() {
    var o = $('#pic-list-items').children().length;
    $('#pic-list-items').append(`
    <tr id="pic-item-` + o + `">
      <td>
        <input name="photos[]" id="pic-` + o + `" type="text" class="form-control" placeholder="Nombre de la foto">
      </td>
      <td>
        <button class="btn btn-danger btn-sm pic-remove" data-id="` + o + `" id="pic-remove-` + o + `">Eliminar</button>
      </td>
    </tr>
    `);
    return false;
  });
  $('#pic-list-items').delegate('button', 'click', function() {
    var id = $(this).data('id');
    var editRemove = $(this).hasClass('pic-edit-remove');
    $('#pic-item-' + id).remove();
    return false;
  });
  $('#' + $('#role').val()).show();
  $('#role').change(function () {
    var r = $(this).val();
    $('.role-detail').hide('fast', function () {
      $('#' + r).show('fast');
    });
  });
  $('div.dataTables_filter input').css('width', '100%');
  $('#user-list').DataTable({
    "dom": "<'form-group full-width'f>rt<'row'<'col-sm'p>>",
    "responsive": true,
    "language": {
      "decimal": "",
      "searchPlaceholder": "Buscar usuario",
      "emptyTable": "No hay datos que mostrar",
      "info": "Mostrando _START_ - _END_ de _TOTAL_ entradas",
      "infoEmpty": "Mostrando 0 - 0 de 0 entradas",
      "infoFiltered": "(filtradas de _MAX_ entradas)",
      "infoPostFix": "",
      "thousands": ".",
      "lengthMenu": "_MENU_",
      "loadingRecords": "Cargando...",
      "processing": "Procesando...",
      "search": "",
      "zeroRecords": "No hay datos que mostrar",
      "paginate": {
        "first": "Primero",
        "last": "Último",
        "next": "Siguiente",
        "previous": "Previo"
      },
      "aria": {
        "sortAscending": ": Activar para ordenar desde columna ascendente",
        "sortDescending": ": Activar para ordenar desde columna descendente"
      }
    }
  });
  $('#file-list').DataTable({
    "dom": "<'form-group full-width'f>rt<'row'<'col-sm'p>>",
    "responsive": true,
    "language": {
      "decimal": "",
      "searchPlaceholder": "Buscar fichero",
      "emptyTable": "No hay datos que mostrar",
      "info": "Mostrando _START_ - _END_ de _TOTAL_ entradas",
      "infoEmpty": "Mostrando 0 - 0 de 0 entradas",
      "infoFiltered": "(filtradas de _MAX_ entradas)",
      "infoPostFix": "",
      "thousands": ".",
      "lengthMenu": "_MENU_",
      "loadingRecords": "Cargando...",
      "processing": "Procesando...",
      "search": "",
      "zeroRecords": "No hay datos que mostrar",
      "paginate": {
        "first": "Primero",
        "last": "Último",
        "next": "Siguiente",
        "previous": "Previo"
      },
      "aria": {
        "sortAscending": ": Activar para ordenar desde columna ascendente",
        "sortDescending": ": Activar para ordenar desde columna descendente"
      }
    }
  });
  $('#log-list').DataTable({
    "bSort": false,
    "dom": "<'form-group full-width'f>rt<'row'<'col-sm'p>>",
    "responsive": true,
    "language": {
      "decimal": "",
      "searchPlaceholder": "Buscar acciones",
      "emptyTable": "No hay datos que mostrar",
      "info": "Mostrando _START_ - _END_ de _TOTAL_ entradas",
      "infoEmpty": "Mostrando 0 - 0 de 0 entradas",
      "infoFiltered": "(filtradas de _MAX_ entradas)",
      "infoPostFix": "",
      "thousands": ".",
      "lengthMenu": "_MENU_",
      "loadingRecords": "Cargando...",
      "processing": "Procesando...",
      "search": "",
      "zeroRecords": "No hay datos que mostrar",
      "paginate": {
        "first": "Primero",
        "last": "Último",
        "next": "Siguiente",
        "previous": "Previo"
      },
      "aria": {
        "sortAscending": ": Activar para ordenar desde columna ascendente",
        "sortDescending": ": Activar para ordenar desde columna descendente"
      }
    }
  });
  $('#box-list').DataTable({
    "dom": "<'form-group full-width'f>rt<'row'<'col-sm'p>>",
    "responsive": true,
    "language": {
      "decimal": "",
      "searchPlaceholder": "Buscar cajas",
      "emptyTable": "No hay datos que mostrar",
      "info": "Mostrando _START_ - _END_ de _TOTAL_ entradas",
      "infoEmpty": "Mostrando 0 - 0 de 0 entradas",
      "infoFiltered": "(filtradas de _MAX_ entradas)",
      "infoPostFix": "",
      "thousands": ".",
      "lengthMenu": "_MENU_",
      "loadingRecords": "Cargando...",
      "processing": "Procesando...",
      "search": "",
      "zeroRecords": "No hay datos que mostrar",
      "paginate": {
        "first": "Primero",
        "last": "Último",
        "next": "Siguiente",
        "previous": "Previo"
      },
      "aria": {
        "sortAscending": ": Activar para ordenar desde columna ascendente",
        "sortDescending": ": Activar para ordenar desde columna descendente"
      }
    }
  });
  $('#pic-list').DataTable({
    "dom": "<'form-group full-width'f>rt<'row'<'col-sm'p>>",
    "responsive": true,
    "language": {
      "decimal": "",
      "searchPlaceholder": "Buscar fotos",
      "emptyTable": "No hay datos que mostrar",
      "info": "Mostrando _START_ - _END_ de _TOTAL_ entradas",
      "infoEmpty": "Mostrando 0 - 0 de 0 entradas",
      "infoFiltered": "(filtradas de _MAX_ entradas)",
      "infoPostFix": "",
      "thousands": ".",
      "lengthMenu": "_MENU_",
      "loadingRecords": "Cargando...",
      "processing": "Procesando...",
      "search": "",
      "zeroRecords": "No hay datos que mostrar",
      "paginate": {
        "first": "Primero",
        "last": "Último",
        "next": "Siguiente",
        "previous": "Previo"
      },
      "aria": {
        "sortAscending": ": Activar para ordenar desde columna ascendente",
        "sortDescending": ": Activar para ordenar desde columna descendente"
      }
    }
  });
  $('#work-list').DataTable({
    "dom": "<'form-group full-width'f>rt<'row'<'col-sm'p>>",
    "responsive": true,
    "language": {
      "decimal": "",
      "searchPlaceholder": "Buscar obra",
      "emptyTable": "No hay datos que mostrar",
      "info": "Mostrando _START_ - _END_ de _TOTAL_ entradas",
      "infoEmpty": "Mostrando 0 - 0 de 0 entradas",
      "infoFiltered": "(filtradas de _MAX_ entradas)",
      "infoPostFix": "",
      "thousands": ".",
      "lengthMenu": "_MENU_",
      "loadingRecords": "Cargando...",
      "processing": "Procesando...",
      "search": "",
      "zeroRecords": "No hay datos que mostrar",
      "paginate": {
        "first": "Primero",
        "last": "Último",
        "next": "Siguiente",
        "previous": "Previo"
      },
      "aria": {
        "sortAscending": ": Activar para ordenar desde columna ascendente",
        "sortDescending": ": Activar para ordenar desde columna descendente"
      }
    }
  });
});
