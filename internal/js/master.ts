import TerrainMap from './modules/terrainMap.ts'
document.addEventListener("DOMContentLoaded", function() {
  console.log('document loaded ready!')

  const newMap = new TerrainMap()
  newMap.init()
});
