/*function setClock(){
  var t = new Date();
  
  // I know its ugly... but make it better and this efficient ;)
  if(t.getHours() > 9){
  	$('.clock .digit:nth-child(1) .pixel').removeClass().addClass('pixel').addClass('_'+(t.getHours()+"").split("")[0]);
		$('.clock .digit:nth-child(2) .pixel').removeClass().addClass('pixel').addClass('_'+(t.getHours()+"").split('')[1])
	} else {
		$('.clock .digit:nth-child(1) .pixel').removeClass().addClass('pixel').addClass('_0');
		$('.clock .digit:nth-child(2) .pixel').removeClass().addClass('pixel').addClass('_'+(t.getHours()+"").split('')[0])
	}

if(t.getMinutes() > 9){$('.clock .digit:nth-child(4) .pixel').removeClass().addClass('pixel').addClass('_'+(t.getMinutes()+"").split("")[0]); $('.clock .digit:nth-child(5) .pixel').removeClass().addClass('pixel').addClass('_'+(t.getMinutes()+"").split('')[1]) } else { $('.clock .digit:nth-child(4) .pixel').removeClass().addClass('pixel').addClass('_0'); $('.clock .digit:nth-child(5) .pixel').removeClass().addClass('pixel').addClass('_'+(t.getMinutes()+"").split('')[0]) }

if(t.getSeconds() > 9){ $('.clock .digit:nth-child(7) .pixel').removeClass().addClass('pixel').addClass('_'+(t.getSeconds()+"").split("")[0]); $('.clock .digit:nth-child(8) .pixel').removeClass().addClass('pixel').addClass('_'+(t.getSeconds()+"").split('')[1]) } else {	$('.clock .digit:nth-child(7) .pixel').removeClass().addClass('pixel').addClass('_0'); $('.clock .digit:nth-child(8) .pixel').removeClass().addClass('pixel').addClass('_'+(t.getSeconds()+"").split('')[0]) }
			}

$(document).ready(function(){
	setClock();
  window.setInterval(function(){ setClock(); },1000);
});
*/
