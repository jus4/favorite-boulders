import proj4Config from "./proj4Config"
const ol =  globalThis.ol

const mapLayers = () => {
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
        projection: proj4Config.projection,
        tileGrid: new ol.tilegrid.WMTS({
          origin: ol.extent.getTopLeft(proj4Config.extent),
          resolutions: proj4Config.resolutions,
          matrixIds: proj4Config.matrixIds
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
        projection: proj4Config.projection,
        tileGrid: new ol.tilegrid.WMTS({
          origin: ol.extent.getTopLeft(proj4Config.extent),
          resolutions: proj4Config.resolutions,
          matrixIds: proj4Config.matrixIds
        }),
        style: 'default',
        wrapX: false
      })
    })

    return {
      maastokartta,
      ortokuva
    }
}

export default mapLayers
