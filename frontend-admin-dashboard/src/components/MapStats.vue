<template>
  <div class="map-stats">
    <!-- زر إظهار/إخفاء الإحصائيات -->
    <button 
      class="stats-toggle" 
      @click="toggleStats"
      :class="{ 'stats-toggle--active': isOpen }"
      :title="isOpen ? 'إخفاء الإحصائيات' : 'إظهار الإحصائيات'"
    >
      <BarChart3 :size="20" />
    </button>

    <!-- لوحة الإحصائيات -->
    <div class="stats-panel" :class="{ 'stats-panel--open': isOpen }">
      <div class="stats-header">
        <span class="stats-title">إحصائيات مباشرة</span>
        <button class="stats-close" @click="toggleStats">
          <X :size="16" />
        </button>
      </div>
      
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-value stat-value--primary">{{ totalEmployees }}</div>
          <div class="stat-label">إجمالي الموظفين</div>
        </div>
        
        <div class="stat-item">
          <div class="stat-value stat-value--success">{{ presentEmployees }}</div>
          <div class="stat-label">حاضرون</div>
        </div>
        
        <div class="stat-item">
          <div class="stat-value stat-value--warning">{{ absentEmployees }}</div>
          <div class="stat-label">غائبون</div>
        </div>
        
        <div class="stat-item">
          <div class="stat-value stat-value--info">{{ totalWorksites }}</div>
          <div class="stat-label">مواقع العمل</div>
        </div>
      </div>

      <div class="stats-progress">
        <div class="progress-bar">
          <div 
            class="progress-fill" 
            :style="{ width: attendanceRate + '%' }"
          ></div>
        </div>
        <div class="progress-label">
          <span>نسبة الحضور</span>
          <span class="progress-value">{{ attendanceRate }}%</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { BarChart3, X } from '@lucide/vue'

const props = defineProps({
  employees: { type: Array, default: () => [] },
  worksites: { type: Array, default: () => [] }
})

const isOpen = ref(false)

const totalEmployees = computed(() => props.employees.length)
const presentEmployees = computed(() => props.employees.filter(e => e.status === 'inside').length)
const absentEmployees = computed(() => props.employees.filter(e => e.status === 'outside').length)
const totalWorksites = computed(() => props.worksites.length)

const attendanceRate = computed(() => {
  if (totalEmployees.value === 0) return 0
  return Math.round((presentEmployees.value / totalEmployees.value) * 100)
})

function toggleStats() {
  isOpen.value = !isOpen.value
}
</script>

<style scoped>
.map-stats {
  position: absolute;
  top: 76px;
  right: 16px;
  z-index: 999;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.stats-toggle {
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

.stats-toggle::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.3) 0%, rgba(76, 175, 80, 0) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.stats-toggle:hover {
  transform: translateY(-2px) scale(1.05);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(76, 175, 80, 0.3);
}

.stats-toggle:hover::before {
  opacity: 1;
}

.stats-toggle--active {
  background: linear-gradient(135deg, #4CAF50 0%, #388E3C 100%);
  box-shadow: 0 8px 32px rgba(76, 175, 80, 0.4), 0 0 0 1px rgba(76, 175, 80, 0.5);
}

.stats-toggle--active::before {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.2) 0%, rgba(255, 255, 255, 0) 100%);
  opacity: 1;
}

[data-theme="dark"] .stats-toggle {
  background: linear-gradient(135deg, #0a0a0a 0%, #1a1a1a 100%);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

[data-theme="dark"] .stats-toggle:hover {
  border-color: rgba(76, 175, 80, 0.5);
}

.stats-panel {
  position: absolute;
  top: 56px;
  right: 0;
  width: 260px;
  background: linear-gradient(135deg, rgba(26, 26, 46, 0.98) 0%, rgba(22, 33, 62, 0.98) 100%);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(255, 255, 255, 0.1);
  padding: 14px;
  opacity: 0;
  visibility: hidden;
  transform: translateY(-10px) scale(0.95);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.stats-panel--open {
  opacity: 1;
  visibility: visible;
  transform: translateY(0) scale(1);
}

[data-theme="dark"] .stats-panel {
  background: linear-gradient(135deg, rgba(10, 10, 10, 0.98) 0%, rgba(26, 26, 26, 0.98) 100%);
  border-color: rgba(255, 255, 255, 0.1);
}

.stats-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

[data-theme="dark"] .stats-header {
  border-color: rgba(255, 255, 255, 0.1);
}

.stats-title {
  font-size: 13px;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: 0.5px;
}

[data-theme="dark"] .stats-title {
  color: #ffffff;
}

.stats-close {
  width: 26px;
  height: 26px;
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

.stats-close:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: rotate(90deg);
}

[data-theme="dark"] .stats-close {
  background: rgba(255, 255, 255, 0.1);
  color: #ffffff;
}

[data-theme="dark"] .stats-close:hover {
  background: rgba(255, 255, 255, 0.2);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
  margin-bottom: 10px;
}

.stat-item {
  text-align: center;
  padding: 10px 8px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.05);
  transition: all 0.3s ease;
}

.stat-item:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
}

[data-theme="dark"] .stat-item {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.05);
}

[data-theme="dark"] .stat-item:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.15);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.3);
}

.stat-value {
  font-size: 20px;
  font-weight: 800;
  line-height: 1;
  margin-bottom: 3px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
}

.stat-value--primary {
  color: #2196F3;
}

.stat-value--success {
  color: #4CAF50;
}

.stat-value--warning {
  color: #FF9800;
}

.stat-value--info {
  color: #00BCD4;
}

[data-theme="dark"] .stat-value--primary {
  color: #2196F3;
}

[data-theme="dark"] .stat-value--success {
  color: #4CAF50;
}

[data-theme="dark"] .stat-value--warning {
  color: #FF9800;
}

[data-theme="dark"] .stat-value--info {
  color: #00BCD4;
}

.stat-label {
  font-size: 10px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.7);
  text-transform: uppercase;
  letter-spacing: 0.6px;
}

[data-theme="dark"] .stat-label {
  color: rgba(255, 255, 255, 0.7);
}

.stats-progress {
  padding-top: 8px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

[data-theme="dark"] .stats-progress {
  border-color: rgba(255, 255, 255, 0.1);
}

.progress-bar {
  height: 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 999px;
  overflow: hidden;
  margin-bottom: 6px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

[data-theme="dark"] .progress-bar {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.05);
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #4CAF50 0%, #8BC34A 100%);
  border-radius: 999px;
  transition: width 0.5s ease;
  box-shadow: 0 0 20px rgba(76, 175, 80, 0.4);
}

[data-theme="dark"] .progress-fill {
  background: linear-gradient(90deg, #4CAF50 0%, #8BC34A 100%);
  box-shadow: 0 0 20px rgba(76, 175, 80, 0.5);
}

.progress-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 10px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.7);
  letter-spacing: 0.5px;
}

[data-theme="dark"] .progress-label {
  color: rgba(255, 255, 255, 0.7);
}

.progress-value {
  font-weight: 800;
  color: #4CAF50;
}

[data-theme="dark"] .progress-value {
  color: #4CAF50;
}

/* Responsive */
@media (max-width: 768px) {
  .map-stats {
    top: 76px;
    right: 8px;
  }
  
  .stats-panel {
    width: 240px;
    right: -8px;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 6px;
  }
  
  .stat-value {
    font-size: 16px;
  }
  
  .stat-label {
    font-size: 8px;
  }
}
</style>
