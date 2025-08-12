const GET_SECTORS_API_URL  = '/api/get-sectors/'

async function getLocations() {
  try {
    const response = await fetch(GET_SECTORS_API_URL)
    if (!response.ok) {
      throw new Error('Fetch error')
    }
    const json = await response.json();
    return json
  } catch(err) {
    console.warn(err)
  }
}

const ol = globalThis.ol

const sectors = async () => {
  const locationsData = await getLocations()
  const locations = locationsData.map((s) => {
    return {
      name: s.name,
      coords: [s.longitude, s.latitude],
      sectorId: s.sector_id,
      cragName: s.crag_name,
      cragId: s.crag_id
    }
  })
  const sectorMapFeatures = locations.map(loc => {
    const transformed = ol.proj.transform(loc.coords, 'EPSG:4326', 'EPSG:3067');
    return new ol.Feature({
      geometry: new ol.geom.Point(transformed),
      name: loc.name,
      sectorId: loc.sectorId,
      cragName: loc.cragName,
      cragId: loc.cragId,
    });
  });

  const vectorSource = new ol.source.Vector({
    features: sectorMapFeatures
  });

  // const vectorLayer = new ol.layer.Vector({
  //   source: vectorSource,
  //   style: new ol.style.Style({
  //     image: new ol.style.Circle({
  //       radius: 5,
  //       fill: new ol.style.Fill({ color: 'blue' }),
  //       stroke: new ol.style.Stroke({ color: 'white', width: 1 })
  //     })
  //   })
  // });

  const vectorLayer = new ol.layer.WebGLVector({
    source: vectorSource,
    style: {
      'circle-radius': 5,
      'circle-fill-color': 'blue',
      'circle-stroke-color': 'white',
      'circle-stroke-width': 1
    }
  });

  return vectorLayer
}

export default sectors()
