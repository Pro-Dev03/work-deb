<template>
  <div 
    class="pull-to-refresh"
    @touchstart="handleTouchStart"
    @touchmove="handleTouchMove"
    @touchend="handleTouchEnd"
    @mousedown="handleTouchStart"
    @mousemove="handleTouchMove"
    @mouseup="handleTouchEnd"
  >
    <div 
      class="pull-indicator"
      :style="{ 
        transform: `translateY(${pullDistance}px)`,
        opacity: Math.min(pullDistance / 60, 1)
      }"
    >
      <div class="pull-icon" :class="{ rotating: isRefreshing }">
        <svg v-if="isRefreshing" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 4v1M12 20v1M4 12h1M19 12h1M6.34 6.34l.71.71M17.24 17.24l.71.71M6.34 17.65l.71-.71M17.24 6.65l.71-.71"/>
        </svg>
      </div>
      <div class="pull-text">{{ refreshText }}</div>
    </div>
    
    <div class="pull-content" :style="{ transform: `translateY(${Math.max(0, pullDistance - 60)}px)` }">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '../services/i18n'

const { t } = useI18n()

const props = defineProps({
  threshold: {
    type: Number,
    default: 80
  },
  refreshText: {
    type: String,
    default: ''
  },
  refreshingText: {
    type: String,
    default: ''
  },
  releaseText: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['refresh'])

const isPulling = ref(false)
const isRefreshing = ref(false)
const pullDistance = ref(0)
const startY = ref(0)
const currentY = ref(0)

const displayText = computed(() => {
  if (isRefreshing.value) return props.refreshingText || t('refreshing')
  if (pullDistance.value >= props.threshold) return props.releaseText || t('release_to_refresh')
  return props.refreshText || t('pull_to_refresh')
})

const refreshText = computed(() => displayText.value)

function handleTouchStart(e) {
  if (isRefreshing.value) return
  
  const y = e.touches ? e.touches[0].clientY : e.clientY
  const scrollTop = e.target?.scrollTop || document.documentElement.scrollTop
  
  // فقط إذا كنا في أعلى الصفحة
  if (scrollTop <= 0) {
    startY.value = y
    isPulling.value = true
    pullDistance.value = 0
  }
}

function handleTouchMove(e) {
  if (!isPulling.value || isRefreshing.value) return
  
  const y = e.touches ? e.touches[0].clientY : e.clientY
  currentY.value = y
  
  const distance = y - startY.value
  
  // حساب المسافة مع تخفيف الحركة
  if (distance > 0) {
    pullDistance.value = Math.min(distance * 0.5, 150)
    
    // منع السحب الافتراضي
    if (e.cancelable) {
      e.preventDefault()
    }
  }
}

function handleTouchEnd() {
  if (!isPulling.value) return
  
  isPulling.value = false
  
  // إذا وصلنا لعتبة التحديث
  if (pullDistance.value >= props.threshold) {
    triggerRefresh()
  }
  
  // إعادة المسافة للصفر
  pullDistance.value = 0
}

async function triggerRefresh() {
  isRefreshing.value = true
  pullDistance.value = 60
  
  try {
    await emit('refresh')
  } catch (error) {
    console.error('Refresh failed:', error)
  } finally {
    setTimeout(() => {
      isRefreshing.value = false
      pullDistance.value = 0
    }, 500)
  }
}
</script>

<style scoped>
.pull-to-refresh {
  position: relative;
  width: 100%;
  min-height: 100vh;
  overflow: hidden;
}

.pull-indicator {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  z-index: 1000;
  background: linear-gradient(180deg, var(--primary-50) 0%, transparent 100%);
}

.pull-icon {
  width: 24px;
  height: 24px;
  color: var(--primary-600);
  margin-bottom: 4px;
  transition: transform 0.3s ease;
}

.pull-icon.rotating {
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.pull-text {
  font-size: var(--text-xs);
  color: var(--primary-600);
  font-weight: var(--font-semibold);
  text-align: center;
}

.pull-content {
  transition: transform 0.3s ease;
  min-height: 100vh;
}

@media (prefers-reduced-motion: reduce){
  .pull-icon.rotating{ animation: none; }
}

@media (max-width: 480px) {
  .pull-text {
    font-size: 11px;
  }
  
  .pull-icon {
    width: 20px;
    height: 20px;
  }
}
</style>
