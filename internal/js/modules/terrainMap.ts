import mapLayers from "./mapLayers"
import proj4Config from "./proj4Config"
import sectors from "./sectors"

class TerrainMap {
  getSectorsUrl: string
  getRoutesBySectorId: string
  routeInfoVisible: boolean
  sectorInfoContainer: HTMLElement | null
  sectorInfo: HTMLElement | null
  sectorInfoContent: HTMLElement | null
  closeSectorRouteInfo: HTMLButtonElement | null
  mapOrtokuvaBtn: HTMLButtonElement | null
  mapMaastokarttaBtn: HTMLButtonElement | null
  mapZoomInBtn: HTMLButtonElement | null
  mapZoomOutBtn: HTMLButtonElement | null
  userLocationBtn: HTMLButtonElement | null
  searchSectorIdBtns: NodeListOf<HTMLButtonElement> | null
  ol: any
  // proj4jConfig: () => {projection: any, extent:any, resolutions: any, matrixIds: any}
  mapLayers: {ortokuva: any, maastokartta: any}
  sectors: any
  view: any
  map: any
  baseLayerGroup: any
  proj4: any
  userLocationStyle: any
  userLocationLayer: any
  sectorElements: {
    popup: HTMLElement | null
    popupContent: HTMLElement | null
    sectorInfo: HTMLElement | null
    sectorRouteList: HTMLElement | null
  }
  sectorOverlay: any

  constructor() {
    this.ol =  globalThis.ol

    this.getSectorsUrl = '/api/get-sectors/'
    this.getRoutesBySectorId = '/api/routes-by-sector/'
    this.routeInfoVisible = false
    this.sectorInfoContainer = document.getElementById('sector-route-info')
    this.sectorInfo = document.getElementById('sector-info')
    this.sectorInfoContent = document.getElementById('sector-route-info-content')
    this.closeSectorRouteInfo = document.getElementById('close-route-info-btn') as HTMLButtonElement | null
    this.mapOrtokuvaBtn = document.getElementById('ortokuva-map-btn') as HTMLButtonElement | null
    this.mapMaastokarttaBtn = document.getElementById('maastokartta-map-btn') as HTMLButtonElement | null
    this.mapZoomInBtn = document.getElementById('map-zoom-in-btn') as HTMLButtonElement | null
    this.mapZoomOutBtn = document.getElementById('map-zoom-out-btn') as HTMLButtonElement | null
    this.userLocationBtn = document.getElementById('user-loc-btn') as HTMLButtonElement | null
    this.searchSectorIdBtns = document.querySelectorAll('.js-search-route-sector-id')

    this.sectorElements = {
      popup: document.getElementById('popup'),
      popupContent:document.getElementById('popup-content'),
      sectorInfo: document.getElementById('sector-info'),
      sectorRouteList: document.querySelector('.sector-route-list') as HTMLElement
    }

    this.sectorOverlay = new this.ol.Overlay({
      element: this.sectorElements.popup,
      autoPan: true,
      autoPanAnimation: {
        duration: 250
      }
    })

    this.mapLayers = mapLayers()
    this.baseLayerGroup = new this.ol.layer.Group({
      title: 'Base Maps',
      layers: [this.mapLayers.maastokartta, this.mapLayers.ortokuva]
    });
    this.view = new this.ol.View({
      projection: proj4Config.projection,
      center: this.ol.extent.getCenter(proj4Config.extent),
      zoom: 0,
      resolutions: proj4Config.resolutions,
      extent: proj4Config.extent
    })
    this.map = new this.ol.Map({
      target: 'map',
      layers: [this.baseLayerGroup],
      view: this.view,
      controls: [],
    });

    // Style for the user location marker
    this.userLocationStyle = new this.ol.style.Style({
      image: new this.ol.style.Circle({
        radius: 10,
        fill: new this.ol.style.Fill({ color: 'red' }),
        stroke: new this.ol.style.Stroke({ color: 'white', width: 2 })
      })
    });

    // Vector layer to show the user's location
    this.userLocationLayer = new this.ol.layer.Vector({
      source: new this.ol.source.Vector()
    });


    this.proj4 = globalThis.proj4
    this.proj4.defs("EPSG:3067", "+proj=utm +zone=35 +ellps=GRS80 +units=m +no_defs");
    this.ol.proj.proj4.register(this.proj4);
    
    if (!this.sectorInfoContent || !this.sectorInfo || !this.sectorInfoContainer || !this.closeSectorRouteInfo ) return
  }

  async getRoutesBySector(id:string, name: string) {
    try {
      const response = await fetch(`${this.getRoutesBySectorId}${id}?name=${name}`)
      if (!response.ok) {
        throw new Error('Fetch error')
      }
      const html = await response.text();
      return html
    } catch(err) {
      console.log(err)
    }
  }

  showRouteInfo(routeId:string) {
    if (!this.sectorInfoContainer || !this.sectorInfo || !this.sectorInfoContent) return
    const routeInfoContent = this.sectorInfo.querySelector(`[data-route-info="${routeId}"]`)
    this.sectorInfoContent.innerHTML = routeInfoContent?.innerHTML || ''
    this.sectorInfoContainer.classList.remove('hidden')
    this.routeInfoVisible = true
  }

  closeRouteInfoBtn() {
    this.sectorInfoContainer?.classList.add('hidden')
  }

  hideRouteInfo() {
    if (!this.routeInfoVisible ) return
    this.sectorInfoContainer?.classList.add('hidden')
  }

  handleRouteInfoClick(e: MouseEvent) {
    const target = e.currentTarget as HTMLButtonElement;
    const {routeId } = target.dataset
    if (!routeId) return

    this.showRouteInfo(routeId)
  }

  initRouteButtonEventListeners() {
    const routeInfoButtons = document.querySelectorAll('.js-route-info-btn')
    this.closeSectorRouteInfo?.addEventListener('click', this.closeRouteInfoBtn.bind(this))
    if (routeInfoButtons && routeInfoButtons.length > 0 ) {
      routeInfoButtons.forEach((infoBtn) => infoBtn.addEventListener('click', this.handleRouteInfoClick.bind(this)))
    }
  }

  cleanRoutButtonEventListers() {
    const routeInfoButtons = document.querySelectorAll('.js-route-info-btn')
    if (routeInfoButtons && routeInfoButtons.length > 0 ) {
      routeInfoButtons.forEach((infoBtn) => infoBtn.removeEventListener('click', this.handleRouteInfoClick.bind(this)))
    }
  }

  async handleSectorClick(evt) {
    const feature = this.map.forEachFeatureAtPixel(evt.pixel, f => f, { hitTolerance: 10 });
    if (feature && this.sectorElements.popupContent && this.sectorElements.sectorRouteList ) {
      const coordinates = feature.getGeometry().getCoordinates();
      const name = feature.get('name');
      const sectorId = feature.get('sectorId');
      const cragName = feature.get('cragName');
      const resultHtml = await this.getRoutesBySector(sectorId, name)
      this.cleanRoutButtonEventListers()
      this.hideRouteInfo()
      this.sectorElements.sectorRouteList.innerHTML = resultHtml || ''
      this.initRouteButtonEventListeners()
      this.sectorElements.popupContent.innerHTML = `<strong>${cragName}<br/>${name}</strong>`;
      this.sectorElements.sectorInfo?.classList.remove('hidden')
      this.sectorOverlay.setPosition(coordinates);
    } else {
      this.sectorElements.sectorInfo?.classList.add('hidden')
      this.sectorOverlay.setPosition(undefined); // Hide popup
    }

  }

  setMapLocation(longitude:number, latitude:number) {
    const transformed = this.ol.proj.transform([longitude, latitude], 'EPSG:4326', 'EPSG:3067');
    this.map.getView().setCenter(transformed)
    this.map.getView().setZoom(8)
  }

  showUserLocation(position) {
    const { latitude, longitude } = position.coords;

    // EPSG:4326 (WGS84) to EPSG:3067 (e.g. Finnish national grid)
    const transformed = this.ol.proj.transform([longitude, latitude], 'EPSG:4326', 'EPSG:3067');

    const userFeature = new this.ol.Feature(new this.ol.geom.Point(transformed));
    userFeature.setStyle(this.userLocationStyle);

    this.userLocationLayer.getSource().clear();
    this.userLocationLayer.getSource().addFeature(userFeature);
    this.setMapLocation(longitude, latitude)

  }

  topMenuEventListeners() {
    if (this.mapOrtokuvaBtn && this.mapMaastokarttaBtn) {
      this.mapOrtokuvaBtn.addEventListener('click', () => {
        this.mapLayers.ortokuva.setVisible(true)
        this.mapLayers.maastokartta.setVisible(false)
      })
      this.mapMaastokarttaBtn.addEventListener('click', () => {
        this.mapLayers.ortokuva.setVisible(false)
        this.mapLayers.maastokartta.setVisible(true)
      })
    }

    // init custom zoom buttons 
    if (this.mapZoomInBtn && this.mapZoomOutBtn) {
      this.mapZoomInBtn.addEventListener('click', () => {
        let zoom = this.view.getZoom();
        this.view.setZoom(zoom + 1);
      })

      this.mapZoomOutBtn.addEventListener('click', () => {
        let zoom = this.view.getZoom();
        this.view.setZoom(zoom - 1);
      })
    }

    if (this.userLocationBtn) {
      this.userLocationBtn.addEventListener('click', () => {
        if ('geolocation' in navigator) {
          navigator.geolocation.getCurrentPosition(this.showUserLocation.bind(this), () => {}, {
            enableHighAccuracy: true
          });
        } else {
          alert('Geolocation is not supported by your browser.');
        }
      })
    }
  }

  zoomToSector(sectorId:number, vectorLayer:any, map:any) {
    const source = vectorLayer.getSource();
    const features = source.getFeatures();
  
    const targetFeature = features.find(f => f.get('sectorId') == sectorId);
  
    if (targetFeature) {
      const geometry = targetFeature.getGeometry();
      const size = map.getSize();
  
      map.getView().fit(geometry, {
        size: size,
        padding: [50, 50, 50, 50], // Optional
        maxZoom: 14,               // Optional: set a maximum zoom level
        duration: 1000             // Optional: animate the zoom
      });
      globalThis?.closeQuickSearchModal()
    } else {
      console.warn(`Sector with ID ${sectorId} not found.`);
    }
  }

  // Route search 
  searchBarRouteSearch(e) {
    const { target } = e
    const button = target.closest('.js-search-route-sector-id');
    if (!button) return

    const { sectorId } = button.dataset
    if (!sectorId) return

    this.zoomToSector(sectorId, this.sectors, this.map)
  }

  async init() {
    if ( typeof globalThis?.proj4 === 'undefined' || typeof globalThis?.ol === 'undefined') {
      console.warn('failed to init map')
    }
    this.sectors = await sectors

    this.topMenuEventListeners()
    this.map.on('singleclick', this.handleSectorClick.bind(this))
    this.map.addOverlay(this.sectorOverlay);
    this.map.getView().fit(proj4Config.extent, { size: this.map.getSize() });
    this.map.addLayer(this.sectors);
    this.map.addLayer(this.userLocationLayer);
    window.addEventListener('click', this.searchBarRouteSearch.bind(this))

  }
}

export default TerrainMap
