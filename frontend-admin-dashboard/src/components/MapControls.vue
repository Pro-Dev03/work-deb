<template>
  <div class="map-controls">
    <!-- زر إظهار/إخفاء التحكم -->
    <button 
      class="controls-toggle" 
      @click="toggleControls"
      :class="{ 'controls-toggle--active': isOpen }"
      :title="isOpen ? 'إخفاء الخيارات' : 'إظهار الخيارات'"
    >
      <Layers :size="20" />
    </button>

    <!-- لوحة التحكم -->
    <div class="controls-panel" :class="{ 'controls-panel--open': isOpen }">
      <div class="controls-header">
        <span class="controls-title">إعدادات الخريطة</span>
        <button class="controls-close" @click="toggleControls">
          <X :size="14" />
        </button>
      </div>

      <!-- تبديل نوع الخريطة -->
      <div class="control-section">
        <label class="control-label">نوع الخريطة</label>
        <div class="map-type-grid">
          <button
            v-for="mapType in mapTypes"
            :key="mapType.id"
            class="map-type-btn"
            :class="{ 'map-type-btn--active': selectedMapType === mapType.id }"
            @click="selectMapType(mapType.id)"
            :title="mapType.name"
          >
            <span class="map-type-icon">{{ mapType.icon }}</span>
            <span class="map-type-name">{{ mapType.name }}</span>
          </button>
        </div>
      </div>

      <!-- الطبقات -->
      <div class="control-section">
        <label class="control-label">الطبقات</label>
        <div class="layers-list">
          <label 
            v-for="layer in layers" 
            :key="layer.id"
            class="layer-item"
          >
            <input 
              type="checkbox" 
              v-model="layer.visible"
              @change="toggleLayer(layer.id)"
              class="layer-checkbox"
            >
            <span class="layer-icon">{{ layer.icon }}</span>
            <span class="layer-name">{{ layer.name }}</span>
          </label>
        </div>
      </div>

      <!-- الإجراءات السريعة -->
      <div class="control-section">
        <label class="control-label">إجراءات سريعة</label>
        <div class="quick-actions">
          <button class="action-btn" @click="centerMap" title="توسيط الخريطة">
            <Compass :size="16" />
            <span>توسيط</span>
          </button>
          <button class="action-btn" @click="fitBounds" title="ملائمة للكل">
            <Maximize2 :size="16" />
            <span>ملائمة</span>
          </button>
          <button class="action-btn" @click="toggleFullscreen" title="ملء الشاشة">
            <Expand :size="16" />
            <span>ملء الشاشة</span>
          </button>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Layers, X, Compass, Maximize2, Expand } from '@lucide/vue'

const props = defineProps({
  currentZoom: { type: Number, default: 7 },
  minZoom: { type: Number, default: 2 },
  maxZoom: { type: Number, default: 18 }
})

const emit = defineEmits(['changeMapType', 'toggleLayer', 'centerMap', 'fitBounds', 'toggleFullscreen'])

const isOpen = ref(false)
const selectedMapType = ref('glass')

const mapTypes = [
  { id: 'osm', name: 'Street', icon: '🗺️' },
  { id: 'minimal', name: 'Minimal', icon: '✨' },
  { id: 'glass', name: 'Glass', icon: '💎' },
  { id: 'satellite', name: 'Satellite', icon: '🛰️' },
  { id: 'dark', name: 'Dark', icon: '🌙' }
]

const layers = ref([
  { id: 'employees', name: 'الموظفين', icon: '👤', visible: true },
  { id: 'worksites', name: 'مواقع العمل', icon: '🏢', visible: true }
])

function toggleControls() {
  isOpen.value = !isOpen.value
}

function selectMapType(typeId) {
  selectedMapType.value = typeId
  emit('changeMapType', typeId)
}

function toggleLayer(layerId) {
  const layer = layers.value.find(l => l.id === layerId)
  emit('toggleLayer', layerId, layer.visible)
}

function centerMap() {
  emit('centerMap')
}

function fitBounds() {
  emit('fitBounds')
}

function toggleFullscreen() {
  emit('toggleFullscreen')
}
</script>

<style scoped>
.map-controls {
  position: absolute;
  top: 16px;
  right: 16px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.controls-toggle {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  border: none;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3), 0 0 0 1px rgba(255, 255, 255, 0.1);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.controls-toggle::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(33, 150, 243, 0.3) 0%, rgba(33, 150, 243, 0) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.controls-toggle:hover {
  transform: translateY(-2px) scale(1.05);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(33, 150, 243, 0.3);
}

.controls-toggle:hover::before {
  opacity: 1;
}

.controls-toggle--active {
  background: linear-gradient(135deg, #2196F3 0%, #1976D2 100%);
  box-shadow: 0 8px 32px rgba(33, 150, 243, 0.4), 0 0 0 1px rgba(33, 150, 243, 0.5);
}

.controls-toggle--active::before {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.2) 0%, rgba(255, 255, 255, 0) 100%);
  opacity: 1;
}

[data-theme="dark"] .controls-toggle {
  background: linear-gradient(135deg, #0a0a0a 0%, #1a1a1a 100%);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

[data-theme="dark"] .controls-toggle:hover {
  border-color: rgba(33, 150, 243, 0.5);
}

.controls-panel {
  position: absolute;
  top: 56px;
  right: 0;
  width: 300px;
  background: linear-gradient(135deg, rgba(26, 26, 46, 0.98) 0%, rgba(22, 33, 62, 0.98) 100%);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(255, 255, 255, 0.1);
  padding: 12px;
  opacity: 0;
  visibility: hidden;
  transform: translateY(-10px) scale(0.95);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.controls-panel--open {
  opacity: 1;
  visibility: visible;
  transform: translateY(0) scale(1);
}

[data-theme="dark"] .controls-panel {
  background: linear-gradient(135deg, rgba(10, 10, 10, 0.98) 0%, rgba(26, 26, 26, 0.98) 100%);
  border-color: rgba(255, 255, 255, 0.1);
}

.controls-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
  padding-bottom: 6px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

[data-theme="dark"] .controls-header {
  border-color: rgba(255, 255, 255, 0.1);
}

.controls-title {
  font-size: 12px;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: 0.5px;
}

[data-theme="dark"] .controls-title {
  color: #ffffff;
}

.controls-close {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: none;
  background: rgba(255, 255, 255, 0.1);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  transition: all 0.2s ease;
}

.controls-close:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: rotate(90deg);
}

[data-theme="dark"] .controls-close {
  background: rgba(255, 255, 255, 0.1);
  color: #ffffff;
}

[data-theme="dark"] .controls-close:hover {
  background: rgba(255, 255, 255, 0.2);
}

.control-section {
  margin-bottom: 6px;
}

.control-section:last-child {
  margin-bottom: 0;
}

.control-label {
  display: block;
  font-size: 10px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.7);
  text-transform: uppercase;
  letter-spacing: 0.6px;
  margin-bottom: 4px;
}

[data-theme="dark"] .control-label {
  color: rgba(255, 255, 255, 0.7);
}

.map-type-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 5px;
}

.map-type-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
  padding: 8px 4px;
  border-radius: 6px;
  border: 2px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.map-type-btn:hover {
  border-color: rgba(33, 150, 243, 0.5);
  background: rgba(33, 150, 243, 0.1);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(33, 150, 243, 0.2);
}

.map-type-btn--active {
  border-color: #2196F3;
  background: linear-gradient(135deg, rgba(33, 150, 243, 0.2) 0%, rgba(33, 150, 243, 0.1) 100%);
  box-shadow: 0 6px 20px rgba(33, 150, 243, 0.3);
}

[data-theme="dark"] .map-type-btn {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
}

[data-theme="dark"] .map-type-btn:hover {
  border-color: rgba(33, 150, 243, 0.5);
  background: rgba(33, 150, 243, 0.15);
}

[data-theme="dark"] .map-type-btn--active {
  border-color: #2196F3;
  background: linear-gradient(135deg, rgba(33, 150, 243, 0.25) 0%, rgba(33, 150, 243, 0.15) 100%);
}

.map-type-icon {
  font-size: 16px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
}

.map-type-name {
  font-size: 9px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
  letter-spacing: 0.3px;
}

[data-theme="dark"] .map-type-name {
  color: rgba(255, 255, 255, 0.9);
}

.layers-list {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.layer-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.layer-item:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.15);
}

[data-theme="dark"] .layer-item:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
}

.layer-checkbox {
  width: 14px;
  height: 14px;
  border-radius: 3px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  cursor: pointer;
  accent-color: #2196F3;
  background: rgba(255, 255, 255, 0.05);
}

[data-theme="dark"] .layer-checkbox {
  border-color: rgba(255, 255, 255, 0.3);
  accent-color: #2196F3;
  background: rgba(255, 255, 255, 0.05);
}

.layer-icon {
  font-size: 12px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
}

.layer-name {
  font-size: 10px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
}

[data-theme="dark"] .layer-name {
  color: rgba(255, 255, 255, 0.9);
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 4px;
}

.action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 6px 4px;
  border-radius: 5px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 8px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
  letter-spacing: 0.3px;
}

.action-btn:hover {
  border-color: rgba(33, 150, 243, 0.5);
  background: rgba(33, 150, 243, 0.15);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(33, 150, 243, 0.2);
}

[data-theme="dark"] .action-btn {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.9);
}

[data-theme="dark"] .action-btn:hover {
  border-color: rgba(33, 150, 243, 0.5);
  background: rgba(33, 150, 243, 0.2);
  box-shadow: 0 6px 20px rgba(33, 150, 243, 0.3);
}

/* Responsive */
@media (max-width: 768px) {
  .controls-panel {
    width: 260px;
    right: -8px;
    padding: 10px;
  }
  
  .controls-header {
    margin-bottom: 6px;
    padding-bottom: 5px;
  }
  
  .controls-title {
    font-size: 11px;
  }
  
  .controls-close {
    width: 22px;
    height: 22px;
  }
  
  .control-section {
    margin-bottom: 5px;
  }
  
  .control-label {
    font-size: 9px;
    margin-bottom: 3px;
  }
  
  .map-type-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 4px;
  }
  
  .map-type-btn {
    padding: 6px 3px;
    gap: 2px;
  }
  
  .map-type-icon {
    font-size: 14px;
  }
  
  .map-type-name {
    font-size: 8px;
  }
  
  .layers-list {
    gap: 2px;
  }
  
  .layer-item {
    padding: 5px 6px;
    gap: 5px;
  }
  
  .layer-checkbox {
    width: 12px;
    height: 12px;
  }
  
  .layer-icon {
    font-size: 11px;
  }
  
  .layer-name {
    font-size: 9px;
  }
  
  .quick-actions {
    gap: 3px;
  }
  
  .action-btn {
    padding: 5px 3px;
    gap: 2px;
    font-size: 7px;
  }
}
</style>
