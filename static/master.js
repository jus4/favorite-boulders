(() => {
  // internal/js/modules/proj4Config.ts
  var proj4jConfig = () => {
    const proj4 = globalThis.proj4;
    const ol3 = globalThis.ol;
    proj4.defs("EPSG:3067", "+proj=utm +zone=35 +ellps=GRS80 +units=m +no_defs");
    ol3.proj.proj4.register(proj4);
    const projection = ol3.proj.get("EPSG:3067");
    const extent = [-548576, 6291456, 1548576, 8388608];
    const resolutions = [
      8192,
      4096,
      2048,
      1024,
      512,
      256,
      128,
      64,
      32,
      16,
      8,
      4,
      2,
      1,
      0.5
    ];
    const matrixIds = resolutions.map((_, i) => i);
    return {
      projection,
      extent,
      resolutions,
      matrixIds
    };
  };
  var proj4Config_default = proj4jConfig();

  // internal/js/modules/mapLayers.ts
  var ol = globalThis.ol;
  var mapLayers = () => {
    const maastokartta = new ol.layer.Tile({
      title: "Maastokartta",
      visible: true,
      type: "base",
      source: new ol.source.WMTS({
        url: "/proxy/wmts/1.0.0/maastokartta/default/ETRS-TM35FIN/{TileMatrix}/{TileRow}/{TileCol}.png",
        requestEncoding: "REST",
        layer: "maastokartta",
        matrixSet: "ETRS-TM35FIN",
        format: "image/png",
        projection: proj4Config_default.projection,
        tileGrid: new ol.tilegrid.WMTS({
          origin: ol.extent.getTopLeft(proj4Config_default.extent),
          resolutions: proj4Config_default.resolutions,
          matrixIds: proj4Config_default.matrixIds
        }),
        style: "default",
        wrapX: false
      })
    });
    const ortokuva = new ol.layer.Tile({
      title: "Ortokuva",
      visible: false,
      type: "base",
      source: new ol.source.WMTS({
        url: "/proxy/wmts/1.0.0/ortokuva/default/ETRS-TM35FIN/{TileMatrix}/{TileRow}/{TileCol}.png",
        requestEncoding: "REST",
        layer: "ortokuva",
        matrixSet: "ETRS-TM35FIN",
        format: "image/png",
        projection: proj4Config_default.projection,
        tileGrid: new ol.tilegrid.WMTS({
          origin: ol.extent.getTopLeft(proj4Config_default.extent),
          resolutions: proj4Config_default.resolutions,
          matrixIds: proj4Config_default.matrixIds
        }),
        style: "default",
        wrapX: false
      })
    });
    return {
      maastokartta,
      ortokuva
    };
  };
  var mapLayers_default = mapLayers;

  // internal/js/modules/sectors.ts
  var GET_SECTORS_API_URL = "/api/get-sectors/";
  async function getLocations() {
    try {
      const response = await fetch(GET_SECTORS_API_URL);
      if (!response.ok) {
        throw new Error("Fetch error");
      }
      const json = await response.json();
      return json;
    } catch (err) {
      console.warn(err);
    }
  }
  var ol2 = globalThis.ol;
  var sectors = async () => {
    const locationsData = await getLocations();
    const locations = locationsData.map((s) => {
      return {
        name: s.name,
        coords: [s.longitude, s.latitude],
        sectorId: s.sector_id,
        cragName: s.crag_name,
        cragId: s.crag_id
      };
    });
    const sectorMapFeatures = locations.map((loc) => {
      const transformed = ol2.proj.transform(loc.coords, "EPSG:4326", "EPSG:3067");
      return new ol2.Feature({
        geometry: new ol2.geom.Point(transformed),
        name: loc.name,
        sectorId: loc.sectorId,
        cragName: loc.cragName,
        cragId: loc.cragId
      });
    });
    const vectorSource = new ol2.source.Vector({
      features: sectorMapFeatures
    });
    const vectorLayer = new ol2.layer.Vector({
      source: vectorSource,
      style: new ol2.style.Style({
        image: new ol2.style.Circle({
          radius: 5,
          fill: new ol2.style.Fill({ color: "blue" }),
          stroke: new ol2.style.Stroke({ color: "white", width: 1 })
        })
      })
    });
    return vectorLayer;
  };
  var sectors_default = sectors();

  // internal/js/modules/terrainMap.ts
  var TerrainMap = class {
    getSectorsUrl;
    getRoutesBySectorId;
    routeInfoVisible;
    sectorInfoContainer;
    sectorInfo;
    sectorInfoContent;
    closeSectorRouteInfo;
    mapOrtokuvaBtn;
    mapMaastokarttaBtn;
    mapZoomInBtn;
    mapZoomOutBtn;
    userLocationBtn;
    ol;
    // proj4jConfig: () => {projection: any, extent:any, resolutions: any, matrixIds: any}
    mapLayers;
    view;
    map;
    baseLayerGroup;
    proj4;
    userLocationStyle;
    userLocationLayer;
    sectorElements;
    sectorOverlay;
    constructor() {
      this.ol = globalThis.ol;
      this.getSectorsUrl = "/api/get-sectors/";
      this.getRoutesBySectorId = "/api/routes-by-sector/";
      this.routeInfoVisible = false;
      this.sectorInfoContainer = document.getElementById("sector-route-info");
      this.sectorInfo = document.getElementById("sector-info");
      this.sectorInfoContent = document.getElementById("sector-route-info-content");
      this.closeSectorRouteInfo = document.getElementById("close-route-info-btn");
      this.mapOrtokuvaBtn = document.getElementById("ortokuva-map-btn");
      this.mapMaastokarttaBtn = document.getElementById("maastokartta-map-btn");
      this.mapZoomInBtn = document.getElementById("map-zoom-in-btn");
      this.mapZoomOutBtn = document.getElementById("map-zoom-out-btn");
      this.userLocationBtn = document.getElementById("user-loc-btn");
      this.sectorElements = {
        popup: document.getElementById("popup"),
        popupContent: document.getElementById("popup-content"),
        sectorInfo: document.getElementById("sector-info"),
        sectorRouteList: document.querySelector(".sector-route-list")
      };
      this.sectorOverlay = new this.ol.Overlay({
        element: this.sectorElements.popup,
        autoPan: true,
        autoPanAnimation: {
          duration: 250
        }
      });
      this.mapLayers = mapLayers_default();
      this.baseLayerGroup = new this.ol.layer.Group({
        title: "Base Maps",
        layers: [this.mapLayers.maastokartta, this.mapLayers.ortokuva]
      });
      this.view = new this.ol.View({
        projection: proj4Config_default.projection,
        center: this.ol.extent.getCenter(proj4Config_default.extent),
        zoom: 0,
        resolutions: proj4Config_default.resolutions,
        extent: proj4Config_default.extent
      });
      this.map = new this.ol.Map({
        target: "map",
        layers: [this.baseLayerGroup],
        view: this.view,
        controls: []
      });
      this.userLocationStyle = new this.ol.style.Style({
        image: new this.ol.style.Circle({
          radius: 10,
          fill: new this.ol.style.Fill({ color: "red" }),
          stroke: new this.ol.style.Stroke({ color: "white", width: 2 })
        })
      });
      this.userLocationLayer = new this.ol.layer.Vector({
        source: new this.ol.source.Vector()
      });
      this.proj4 = globalThis.proj4;
      this.proj4.defs("EPSG:3067", "+proj=utm +zone=35 +ellps=GRS80 +units=m +no_defs");
      this.ol.proj.proj4.register(this.proj4);
      if (!this.sectorInfoContent || !this.sectorInfo || !this.sectorInfoContainer || !this.closeSectorRouteInfo) return;
    }
    async getRoutesBySector(id, name) {
      try {
        const response = await fetch(`${this.getRoutesBySectorId}${id}?name=${name}`);
        if (!response.ok) {
          throw new Error("Fetch error");
        }
        const html = await response.text();
        return html;
      } catch (err) {
        console.log(err);
      }
    }
    showRouteInfo(routeId) {
      if (!this.sectorInfoContainer || !this.sectorInfo || !this.sectorInfoContent) return;
      const routeInfoContent = this.sectorInfo.querySelector(`[data-route-info="${routeId}"]`);
      this.sectorInfoContent.innerHTML = routeInfoContent?.innerHTML || "";
      this.sectorInfoContainer.classList.remove("hidden");
      this.routeInfoVisible = true;
    }
    closeRouteInfoBtn() {
      this.sectorInfoContainer?.classList.add("hidden");
    }
    hideRouteInfo() {
      if (!this.routeInfoVisible) return;
      this.sectorInfoContainer?.classList.add("hidden");
    }
    handleRouteInfoClick(e) {
      const target = e.currentTarget;
      const { routeId } = target.dataset;
      if (!routeId) return;
      this.showRouteInfo(routeId);
    }
    initRouteButtonEventListeners() {
      const routeInfoButtons = document.querySelectorAll(".js-route-info-btn");
      this.closeSectorRouteInfo?.addEventListener("click", this.closeRouteInfoBtn.bind(this));
      if (routeInfoButtons && routeInfoButtons.length > 0) {
        routeInfoButtons.forEach((infoBtn) => infoBtn.addEventListener("click", this.handleRouteInfoClick.bind(this)));
      }
    }
    cleanRoutButtonEventListers() {
      const routeInfoButtons = document.querySelectorAll(".js-route-info-btn");
      if (routeInfoButtons && routeInfoButtons.length > 0) {
        routeInfoButtons.forEach((infoBtn) => infoBtn.removeEventListener("click", this.handleRouteInfoClick.bind(this)));
      }
    }
    async handleSectorClick(evt) {
      const feature = this.map.forEachFeatureAtPixel(evt.pixel, (f) => f, { hitTolerance: 10 });
      console.log(this.sectorElements);
      if (feature && this.sectorElements.popupContent && this.sectorElements.sectorRouteList) {
        const coordinates = feature.getGeometry().getCoordinates();
        const name = feature.get("name");
        const sectorId = feature.get("sectorId");
        const cragName = feature.get("cragName");
        const resultHtml = await this.getRoutesBySector(sectorId, name);
        this.cleanRoutButtonEventListers();
        this.hideRouteInfo();
        this.sectorElements.sectorRouteList.innerHTML = resultHtml || "";
        this.initRouteButtonEventListeners();
        this.sectorElements.popupContent.innerHTML = `<strong>${cragName}<br/>${name}</strong>`;
        this.sectorElements.sectorInfo?.classList.remove("hidden");
        this.sectorOverlay.setPosition(coordinates);
      } else {
        this.sectorElements.sectorInfo?.classList.add("hidden");
        this.sectorOverlay.setPosition(void 0);
      }
    }
    setMapLocation(longitude, latitude) {
      const transformed = this.ol.proj.transform([longitude, latitude], "EPSG:4326", "EPSG:3067");
      this.map.getView().setCenter(transformed);
      this.map.getView().setZoom(8);
    }
    showUserLocation(position) {
      const { latitude, longitude } = position.coords;
      const transformed = this.ol.proj.transform([longitude, latitude], "EPSG:4326", "EPSG:3067");
      const userFeature = new this.ol.Feature(new this.ol.geom.Point(transformed));
      userFeature.setStyle(this.userLocationStyle);
      this.userLocationLayer.getSource().clear();
      this.userLocationLayer.getSource().addFeature(userFeature);
      this.setMapLocation(longitude, latitude);
    }
    topMenuEventListeners() {
      if (this.mapOrtokuvaBtn && this.mapMaastokarttaBtn) {
        this.mapOrtokuvaBtn.addEventListener("click", () => {
          this.mapLayers.ortokuva.setVisible(true);
          this.mapLayers.maastokartta.setVisible(false);
        });
        this.mapMaastokarttaBtn.addEventListener("click", () => {
          this.mapLayers.ortokuva.setVisible(false);
          this.mapLayers.maastokartta.setVisible(true);
        });
      }
      if (this.mapZoomInBtn && this.mapZoomOutBtn) {
        this.mapZoomInBtn.addEventListener("click", () => {
          let zoom = this.view.getZoom();
          this.view.setZoom(zoom + 1);
        });
        this.mapZoomOutBtn.addEventListener("click", () => {
          let zoom = this.view.getZoom();
          this.view.setZoom(zoom - 1);
        });
      }
      if (this.userLocationBtn) {
        this.userLocationBtn.addEventListener("click", () => {
          if ("geolocation" in navigator) {
            navigator.geolocation.getCurrentPosition(this.showUserLocation.bind(this), () => {
            }, {
              enableHighAccuracy: true
            });
          } else {
            alert("Geolocation is not supported by your browser.");
          }
        });
      }
    }
    async init() {
      if (typeof globalThis?.proj4 === "undefined" || typeof globalThis?.ol === "undefined") {
        console.warn("failed to init map");
      }
      this.topMenuEventListeners();
      this.map.on("singleclick", this.handleSectorClick.bind(this));
      const sectorLayer = await sectors_default;
      this.map.addOverlay(this.sectorOverlay);
      this.map.getView().fit(proj4Config_default.extent, { size: this.map.getSize() });
      this.map.addLayer(sectorLayer);
      this.map.addLayer(this.userLocationLayer);
    }
  };
  var terrainMap_default = TerrainMap;

  // internal/js/master.ts
  document.addEventListener("DOMContentLoaded", function() {
    const newMap = new terrainMap_default();
    newMap.init();
  });
})();
