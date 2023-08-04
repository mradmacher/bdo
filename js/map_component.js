export class MapComponent {
  constructor(apiKey) {
    (g=>{var h,a,k,p="The Google Maps JavaScript API",c="google",l="importLibrary",q="__ib__",m=document,b=window;b=b[c]||(b[c]={});var d=b.maps||(b.maps={}),r=new Set,e=new URLSearchParams,u=()=>h||(h=new Promise(async(f,n)=>{await (a=m.createElement("script"));e.set("libraries",[...r]+"");for(k in g)e.set(k.replace(/[A-Z]/g,t=>"_"+t[0].toLowerCase()),g[k]);e.set("callback",c+".maps."+q);a.src=`https://maps.${c}apis.com/maps/api/js?`+e;d[q]=f;a.onerror=()=>h=n(Error(p+" could not load."));a.nonce=m.querySelector("script[nonce]")?.nonce||"";m.head.append(a)}));d[l]?console.warn(p+" only loads once. Ignoring:",g):d[l]=(f,...n)=>r.add(f)&&u().then(()=>d[l](f,...n))})({
      key: apiKey,
      v: "weekly",
      // Use the 'v' parameter to indicate the version to use (weekly, beta, alpha, etc.).
      // Add other bootstrap parameters as needed, using camel case.
    });

    this.markers = []
    this.mapCenter = { lat: 52.24, lng: 21.00 }
    this.mapZoom = 6
  }

  clear() {
    this.map.setZoom(this.mapZoom)
    this.map.setCenter(this.mapCenter)
    while (this.markers.length) {
      this.markers.pop().setMap(null)
    }
  }

  addInstallation(installation) {
    const latLng = new google.maps.LatLng(
      installation.Address.Lat,
      installation.Address.Lng
    )
    let marker = new google.maps.Marker({
      map: this.map,
      position: latLng,
      label: installation.Name,
      title: `${installation.Name}\n${installation.Address.Line1}\n${installation.Address.Line2}`,
    })
    const infoWindow = new google.maps.InfoWindow({
      content: `
        <h3>${installation.Name}</h3>
        <p>${installation.Address.Line1}<br />${installation.Address.Line2}</p>
      `,
      ariaLabel: installation.Name,
    })
    marker.addListener("click", () => {
      infoWindow.open({
        anchor: marker,
        map: this.map,
      })
    })
    this.markers.push(marker)
  }

  async initMap(elementId) {
    const { Map } = await google.maps.importLibrary("maps");

    this.map = new Map(document.getElementById(elementId), {
      center: this.mapCenter,
      zoom: this.mapZoom,
    });
  }
}
