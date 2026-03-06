
  var newEntity = function (url, obj, callback) {
    document.getElementById('loadingOverlay').style.display = 'flex';
    setData(url, obj, callback, "POST");
  }
  var setData = function (url, obj, callback, action) {
    
    var response = null;
    const xhr = new XMLHttpRequest()
    //open a get request with the remote server URL
    switch (action) {
      case "POST":
          xhr.open("POST", url);
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
    xhr.onload = function() {
      document.getElementById('loadingOverlay').style.display = 'none';
      if (xhr.status === 200) {
        document.getElementById('successOverlay').style.display = 'flex';
        //parse JSON datax`x 
        if(xhr.responseText.length > 0){
          console.log(xhr.responseText);
          response = JSON.parse(xhr.responseText);
        }
      }
      else if (xhr.status === 400) {
        document.getElementById('failOverlay').style.display = 'flex';
        setTimeout(() => {
          document.getElementById('failOverlay').style.display = 'none';
        }, 2000);
        if(xhr.responseText.length > 0){
          response = JSON.parse(xhr.responseText);
        }
      }else if (xhr.status === 404) {
        console.log("No records found");
      }else{
        document.getElementById('failOverlay').style.display = 'flex';
        setTimeout(() => {
          document.getElementById('failOverlay').style.display = 'none';
        }, 2000);
      }
      //alert(xhr.responseText);
      callback(xhr.status, response);
    };

    //triggered when a network-level error occurs with the request
    xhr.onerror = function() {
      console.log("Network error occurred")
    };

    //triggered periodically as the client receives data
    //used to monitor the progress of the request
    xhr.onprogress = function(e) {
      if (e.lengthComputable) {
        console.log(`${e.loaded} B of ${e.total} B loaded!`)
      } else {
        console.log(`${e.loaded} B loaded!`)
      }
    };
  };
