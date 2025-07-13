(() => {
  // internal/js/modules/terrainMap.ts
  var TerrainMap = class {
    getSectorsUrl;
    getRoutesBySectorId;
    routeInfoVisible;
    sectorInfoContainer;
    sectorInfo;
    sectorInfoContent;
    closeSectorRouteInfo;
    constructor() {
      this.getSectorsUrl = "/api/get-sectors/";
      this.getRoutesBySectorId = "/api/routes-by-sector/";
      this.routeInfoVisible = false;
      this.sectorInfoContainer = document.getElementById("sector-route-info");
      this.sectorInfo = document.getElementById("sector-info");
      this.sectorInfoContent = document.getElementById("sector-route-info-content");
      this.closeSectorRouteInfo = document.getElementById("close-route-info-btn");
      if (!this.sectorInfoContent || !this.sectorInfo || !this.sectorInfoContainer || !this.closeSectorRouteInfo) return;
    }
    async getLocations() {
      try {
        const response = await fetch(this.getSectorsUrl);
        if (!response.ok) {
          throw new Error("Fetch error");
        }
        const json = await response.json();
        return json;
      } catch (err) {
        console.log(err);
      }
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
    async init() {
      if (typeof globalThis?.proj4 === "undefined" || typeof globalThis?.ol === "undefined") {
        console.warn("failed to init map");
      }
      const proj4 = globalThis.proj4;
      const ol = globalThis.ol;
      proj4.defs("EPSG:3067", "+proj=utm +zone=35 +ellps=GRS80 +units=m +no_defs");
      ol.proj.proj4.register(proj4);
      const projection = ol.proj.get("EPSG:3067");
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
      const locationsData = await this.getLocations();
      const locations = locationsData.map((s) => {
        return {
          name: s.name,
          coords: [s.longitude, s.latitude],
          sectorId: s.sector_id,
          cragName: s.crag_name
        };
      });
      const features = locations.map((loc) => {
        const transformed = ol.proj.transform(loc.coords, "EPSG:4326", "EPSG:3067");
        return new ol.Feature({
          geometry: new ol.geom.Point(transformed),
          name: loc.name,
          sectorId: loc.sectorId,
          cragName: loc.cragName
        });
      });
      const vectorSource = new ol.source.Vector({
        features
      });
      const vectorLayer = new ol.layer.Vector({
        source: vectorSource,
        style: new ol.style.Style({
          image: new ol.style.Circle({
            radius: 5,
            fill: new ol.style.Fill({ color: "blue" }),
            stroke: new ol.style.Stroke({ color: "white", width: 1 })
          })
        })
      });
      const popup = document.getElementById("popup");
      const popupContent = document.getElementById("popup-content");
      const sectorInfo = document.getElementById("sector-info");
      const sectorRouteList = sectorInfo?.querySelector(".sector-route-list");
      const overlay = new ol.Overlay({
        element: popup,
        autoPan: true,
        autoPanAnimation: {
          duration: 250
        }
      });
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
          projection,
          tileGrid: new ol.tilegrid.WMTS({
            origin: ol.extent.getTopLeft(extent),
            resolutions,
            matrixIds
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
          projection,
          tileGrid: new ol.tilegrid.WMTS({
            origin: ol.extent.getTopLeft(extent),
            resolutions,
            matrixIds
          }),
          style: "default",
          wrapX: false
        })
      });
      const baseLayerGroup = new ol.layer.Group({
        title: "Base Maps",
        layers: [maastokartta, ortokuva]
      });
      const map = new ol.Map({
        target: "map",
        layers: [baseLayerGroup],
        view: new ol.View({
          projection,
          center: ol.extent.getCenter(extent),
          zoom: 0,
          resolutions,
          extent
        })
      });
      map.on("singleclick", async (evt) => {
        const feature = map.forEachFeatureAtPixel(evt.pixel, (f) => f, { hitTolerance: 10 });
        if (feature && popupContent && sectorRouteList) {
          const coordinates = feature.getGeometry().getCoordinates();
          const name = feature.get("name");
          const sectorId = feature.get("sectorId");
          const cragName = feature.get("cragName");
          const resultHtml = await this.getRoutesBySector(sectorId, name);
          this.cleanRoutButtonEventListers();
          this.hideRouteInfo();
          sectorRouteList.innerHTML = resultHtml || "";
          this.initRouteButtonEventListeners();
          popupContent.innerHTML = `<strong>${cragName}<br/>${name}</strong>`;
          sectorInfo?.classList.remove("hidden");
          overlay.setPosition(coordinates);
        } else {
          sectorInfo?.classList.add("hidden");
          overlay.setPosition(void 0);
        }
      });
      const layerSwitcher = new ol.control.LayerSwitcher({
        tipLabel: "Layers",
        groupSelectStyle: "children"
      });
      map.addControl(layerSwitcher);
      map.addOverlay(overlay);
      map.getView().fit(extent, { size: map.getSize() });
      map.addLayer(vectorLayer);
    }
  };
  var terrainMap_default = TerrainMap;

  // internal/js/master.ts
  document.addEventListener("DOMContentLoaded", function() {
    console.log("document loaded ready!");
    const newMap = new terrainMap_default();
    newMap.init();
  });
})();
