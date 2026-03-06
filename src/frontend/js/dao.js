var getData = function (url, callback, hideDialogs) {
  if(!hideDialogs){
    document.getElementById('loadingOverlay').style.display = 'flex';
  }
  
  var response = null;
  const xhr = new XMLHttpRequest()
  //open a get request with the remote server URL
  xhr.open("GET", url);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.setRequestHeader("Accept", "application/json");

  //send the Http request
  xhr.send();

  //EVENT HANDLERS

  //triggered when the response is completed
  xhr.onload = function () {
    if(!hideDialogs){
      document.getElementById('loadingOverlay').style.display = 'none';
    }
    if(xhr.responseURL.includes("login")){
      alert("Parece que la sesión se ha vencido. Por favor ingresa nuevamente");
      location.assign("/");
    }
    if (xhr.status === 200) {
      if(!hideDialogs){
        document.getElementById('successOverlay').style.display = 'flex';
      }
      //parse JSON datax`x 
      if(xhr.responseText.length > 0){
        response = JSON.parse(xhr.responseText);
      }

    } else if (xhr.status === 404) {
      console.log("No records found")
    }else{
      if(!hideDialogs){
        
        document.getElementById('failOverlay').style.display = 'flex';
        setTimeout(() => {
          document.getElementById('failOverlay').style.display = 'none';
          window.scrollTo(0, 0);
        }, 2000);
      }
    }
    callback(xhr.status, response);
  };

  //triggered when a network-level error occurs with the request
  xhr.onerror = function () {
    console.log("Network error occurred")
  };

  //triggered periodically as the client receives data
  //used to monitor the progress of the request
  xhr.onprogress = function (e) {
    if (e.lengthComputable) {
      console.log(`${e.loaded} B of ${e.total} B loaded!`)
    } else {
      console.log(`${e.loaded} B loaded!`)
    }
  };
};
var getEntity = function (url, callback, hideDialogs) {
  getData(url, callback, hideDialogs)
}
var removeEntity = function (url, obj, callback) {
  setData(url, obj, callback, "DELETE");
}
var updateEntity = function (url, obj, callback) {
  setData(url, obj, callback, "PUT");
}
var sendMultipartForm = function (url, obj, callback) {
  setUploadData(url, obj, callback, "POST");
}
var newEntity = function (url, obj, callback) {
  setData(url, obj, callback, "POST");
}
var setData = function (url, obj, callback, action, logout) {
  
  document.getElementById('loadingOverlay').style.display = 'flex';
  var response = null;
  const xhr = new XMLHttpRequest()
  //open a get request with the remote server URL
  switch (action) {
    case "POST":
      xhr.open("POST", url);
      break;
    case "PUT":
      xhr.open("PUT", url);
      break;
    case "DELETE":
      xhr.open("DELETE", url);
      break;
    default:
      break;
  }


  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.setRequestHeader("Accept", "application/json");
  

  //send the Http request
  //console.log(JSON.stringify(obj))

  xhr.send(JSON.stringify(obj));


  //EVENT HANDLERS

  //triggered when the response is completed
  xhr.onload = function () {
    document.getElementById('loadingOverlay').style.display = 'none';
    if(!logout && xhr.responseURL.includes("login")){
      alert("Parece que la sesión se ha vencido. Por favor ingresa nuevamente");
      location.assign("/");
    }
    
    if (xhr.status === 200) {
      document.getElementById('successOverlay').style.display = 'flex';
      //parse JSON datax`x 
      if(xhr.responseText.length > 0){
        response = JSON.parse(xhr.responseText)
      }
    }
    else if (xhr.status === 400) {
      
      document.getElementById('failOverlay').style.display = 'flex';
        setTimeout(() => {
          document.getElementById('failOverlay').style.display = 'none';
          window.scrollTo(0, 0);
        }, 2000);
        if(xhr.responseText.length > 0){
          response = JSON.parse(xhr.responseText)
        }
    } else if (xhr.status === 404) {
      console.log("No records found")
    }else{
      
      document.getElementById('failOverlay').style.display = 'flex';
        setTimeout(() => {
          document.getElementById('failOverlay').style.display = 'none';
          window.scrollTo(0, 0);
        }, 2000);
    }
    //alert(xhr.responseText);
    callback(xhr.status, response);
  };

  //triggered when a network-level error occurs with the request
  xhr.onerror = function () {
    console.log("Network error occurred")
    document.getElementById('loadingOverlay').style.display = 'none';
    callback(xhr.status, response);
  };

  //triggered periodically as the client receives data
  //used to monitor the progress of the request
  xhr.onprogress = function (e) {
    if (e.lengthComputable) {
      console.log(`${e.loaded} B of ${e.total} B loaded!`)
    } else {
      console.log(`${e.loaded} B loaded!`)
    }
  };
};

var setUploadData = function (url, obj, callback, action, logout) {
  document.getElementById('loadingOverlay').style.display = 'flex';
  var response = null;
  const xhr = new XMLHttpRequest()
  
  xhr.open("POST", url);
  
  
  //xhr.setRequestHeader("Content-Type", "application/json");
  //const boundary = '----------FormDataBoundary' + Math.random().toString(16).substr(2);
  //xhr.setRequestHeader("Content-Type", 'multipart/form-data; boundary=' + boundary);
  //xhr.setRequestHeader("Accept", "application/json");
  

  


  //send the Http request
  //console.log(JSON.stringify(obj))

  xhr.send(obj);


  //EVENT HANDLERS

  //triggered when the response is completed
  xhr.onload = function () {
    document.getElementById('loadingOverlay').style.display = 'none';
    if(!logout && xhr.responseURL.includes("login")){
      alert("Parece que la sesión se ha vencido. Por favor ingresa nuevamente");
      location.assign("/");
    }
    
    if (xhr.status === 200) {
      document.getElementById('successOverlay').style.display = 'flex';
      //parse JSON datax`x 
      if(xhr.responseText.length > 0){
        response = JSON.parse(xhr.responseText)
      }
    }
    else if (xhr.status === 400) {
      
      document.getElementById('failOverlay').style.display = 'flex';
        setTimeout(() => {
          document.getElementById('failOverlay').style.display = 'none';
          window.scrollTo(0, 0);
        }, 2000);
        if(xhr.responseText.length > 0){
          response = JSON.parse(xhr.responseText)
        }
    } else if (xhr.status === 404) {
      console.log("No records found")
    }else{
      
      document.getElementById('failOverlay').style.display = 'flex';
        setTimeout(() => {
          document.getElementById('failOverlay').style.display = 'none';
          window.scrollTo(0, 0);
        }, 2000);
    }
    //alert(xhr.responseText);
    callback(xhr.status, response);
  };

  //triggered when a network-level error occurs with the request
  xhr.onerror = function () {
    console.log("Network error occurred")
  };

  //triggered periodically as the client receives data
  //used to monitor the progress of the request
  xhr.onprogress = function (e) {
    if (e.lengthComputable) {
      console.log(`${e.loaded} B of ${e.total} B loaded!`)
    } else {
      console.log(`${e.loaded} B loaded!`)
    }
  };
};

var logout = function (url){
  setData(url,{},function(status,response){
    location.assign("/");
  }, "POST", true)
}
function buildFormDataBody(formData, boundary) {
  let body = '';

  for (const [name, value] of formData) {
    body += '--' + boundary + '\r\n';
    body += 'Content-Disposition: form-data; name="' + name + '"\r\n\r\n';
    body += value + '\r\n';
  }

  body += '--' + boundary + '--\r\n';
alert(body);
  return body;
}

function mostrarContenido(tab) {
  var tabs = document.querySelectorAll('.tabpopup');

  tabs.forEach(function (el) {
    el.classList.remove('active');
  });

  document.getElementById('info-tab').style.display = (tab === 'info') ? 'block' : 'none';
  document.getElementById('rutas-tab').style.display = (tab === 'rutas') ? 'block' : 'none';

  if (tab === 'info') {
    tabs[0].classList.add('active');
  } else if (tab === 'rutas') {
    tabs[1].classList.add('active');
  }
}