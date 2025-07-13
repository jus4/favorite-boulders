class TerrainMap {
  getSectorsUrl: string
  getRoutesBySectorId: string
  routeInfoVisible: boolean
  sectorInfoContainer: HTMLElement | null
  sectorInfo: HTMLElement | null
  sectorInfoContent: HTMLElement | null
  closeSectorRouteInfo: HTMLButtonElement | null

  constructor() {
    this.getSectorsUrl = '/api/get-sectors'
    this.getRoutesBySectorId = '/api/routes-by-sector/'
    this.routeInfoVisible = false
    this.sectorInfoContainer = document.getElementById('sector-route-info')
    this.sectorInfo = document.getElementById('sector-info')
    this.sectorInfoContent = document.getElementById('sector-route-info-content')
    this.closeSectorRouteInfo = document.getElementById('close-route-info-btn') as HTMLButtonElement | null
    
    if (!this.sectorInfoContent || !this.sectorInfo || !this.sectorInfoContainer || !this.closeSectorRouteInfo ) return
  }

  async getLocations() {
    try {
      const response = await fetch(this.getSectorsUrl)
      if (!response.ok) {
        throw new Error('Fetch error')
      }
      const json = await response.json();
      return json
    } catch(err) {
      console.log(err)
    }
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

  async init() {
    if ( typeof globalThis?.proj4 === 'undefined' || typeof globalThis?.ol === 'undefined') {
      console.warn('failed to init map')
    }
    const proj4 = globalThis.proj4
    const ol = globalThis.ol

    proj4.defs("EPSG:3067", "+proj=utm +zone=35 +ellps=GRS80 +units=m +no_defs");
    ol.proj.proj4.register(proj4);

    const projection = ol.proj.get('EPSG:3067');
    const extent = [-548576.000000, 6291456.000000, 1548576.000000, 8388608.000000];

    const resolutions = [
      8192, 4096, 2048, 1024, 512, 256, 128, 64, 32, 16, 8, 4, 2, 1, 0.5
    ];
    const matrixIds = resolutions.map((_, i) => i);

    // TODO do this in the backend please
    const locationsData = await this.getLocations()
    const locations = locationsData.map((s) => {
      return {
        name: s.name,
        coords: [s.longitude, s.latitude],
        sectorId: s.sector_id,
        cragName: s.crag_name,
      }
    })

    const features = locations.map(loc => {
      const transformed = ol.proj.transform(loc.coords, 'EPSG:4326', 'EPSG:3067');
      return new ol.Feature({
        geometry: new ol.geom.Point(transformed),
        name: loc.name,
        sectorId: loc.sectorId,
        cragName: loc.cragName,
      });
    });

    const vectorSource = new ol.source.Vector({
      features: features
    });

    const vectorLayer = new ol.layer.Vector({
      source: vectorSource,
      style: new ol.style.Style({
        image: new ol.style.Circle({
          radius: 7,
          fill: new ol.style.Fill({ color: 'blue' }),
          stroke: new ol.style.Stroke({ color: 'white', width: 1 })
        })
      })
    });

    const popup = document.getElementById('popup');
    const popupContent = document.getElementById('popup-content');
    const sectorInfo = document.getElementById('sector-info')
    const sectorRouteList = sectorInfo?.querySelector('.sector-route-list')

    const overlay = new ol.Overlay({
      element: popup,
      autoPan: true,
      autoPanAnimation: {
        duration: 250
      }
    });

    const maastokartta = new ol.layer.Tile({
      title: 'Maastokartta',
      visible: true,
      type: 'base',
      source: new ol.source.WMTS({
        url: '/proxy/wmts/1.0.0/maastokartta/default/ETRS-TM35FIN/{TileMatrix}/{TileRow}/{TileCol}.png',
        requestEncoding: 'REST',
        layer: "maastokartta",
        matrixSet: 'ETRS-TM35FIN',
        format: 'image/png',
        projection: projection,
        tileGrid: new ol.tilegrid.WMTS({
          origin: ol.extent.getTopLeft(extent),
          resolutions: resolutions,
          matrixIds: matrixIds
        }),
        style: 'default',
        wrapX: false
      })
    })

    const ortokuva = new ol.layer.Tile({
      title: 'Ortokuva',
      visible: false,
      type: 'base',
      source: new ol.source.WMTS({
        url: '/proxy/wmts/1.0.0/ortokuva/default/ETRS-TM35FIN/{TileMatrix}/{TileRow}/{TileCol}.png',
        requestEncoding: 'REST',
        layer: "ortokuva",
        matrixSet: 'ETRS-TM35FIN',
        format: 'image/png',
        projection: projection,
        tileGrid: new ol.tilegrid.WMTS({
          origin: ol.extent.getTopLeft(extent),
          resolutions: resolutions,
          matrixIds: matrixIds
        }),
        style: 'default',
        wrapX: false
      })
    })

    // Group base layers so layer switcher can treat them as a toggle set
    const baseLayerGroup = new ol.layer.Group({
      title: 'Base Maps',
      layers: [maastokartta, ortokuva]
    });


    const map = new ol.Map({
      target: 'map',
      layers: [baseLayerGroup],
      view: new ol.View({
        projection: projection,
        center: ol.extent.getCenter(extent),
        zoom: 0,
        resolutions: resolutions,
        extent: extent
      })
    });

    // Show popup on click
    map.on('singleclick', async (evt) => {
      const feature = map.forEachFeatureAtPixel(evt.pixel, f => f);
      if (feature && popupContent && sectorRouteList ) {
        const coordinates = feature.getGeometry().getCoordinates();
        const name = feature.get('name');
        const sectorId = feature.get('sectorId');
        const cragName = feature.get('cragName');
        const resultHtml = await this.getRoutesBySector(sectorId, name)
        this.cleanRoutButtonEventListers()
        this.hideRouteInfo()
        sectorRouteList.innerHTML = resultHtml || ''
        this.initRouteButtonEventListeners()
        popupContent.innerHTML = `<strong>${cragName}<br/>${name}</strong>`;
        sectorInfo?.classList.remove('hidden')
        overlay.setPosition(coordinates);
      } else {
        sectorInfo?.classList.add('hidden')
        overlay.setPosition(undefined); // Hide popup
      }
    });

    // Add layer switcher control
    const layerSwitcher = new ol.control.LayerSwitcher({
      tipLabel: 'Layers', // Tooltip
      groupSelectStyle: 'children' // Show individual base layers
    });
    map.addControl(layerSwitcher);
    map.addOverlay(overlay);
    map.getView().fit(extent, { size: map.getSize() });
    map.addLayer(vectorLayer);

  }
}

export default TerrainMap
