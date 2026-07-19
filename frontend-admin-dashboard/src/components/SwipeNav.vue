<template>
  <div 
    class="swipe-nav"
    @touchstart="handleTouchStart"
    @touchmove="handleTouchMove"
    @touchend="handleTouchEnd"
    @mousedown="handleTouchStart"
    @mousemove="handleTouchMove"
    @mouseup="handleTouchEnd"
    @mouseleave="handleTouchEnd"
  >
    <div 
      class="swipe-content"
      :style="{ 
        transform: `translateX(${translateX}px)`,
        transition: isDragging ? 'none' : 'transform 0.3s ease'
      }"
    >
      <slot></slot>
    </div>
    
    <div class="swipe-indicators" v-if="showIndicators && totalItems > 1">
      <div 
        v-for="(_, index) in totalItems" 
        :key="index"
        class="indicator"
        :class="{ active: currentIndex === index }"
        @click="goToSlide(index)"
      ></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  itemsPerView: {
    type: Number,
    default: 1
  },
  showIndicators: {
    type: Boolean,
    default: true
  },
  autoPlay: {
    type: Boolean,
    default: false
  },
  autoPlayInterval: {
    type: Number,
    default: 3000
  }
})

const emit = defineEmits(['slide-change', 'swipe-left', 'swipe-right'])

const translateX = ref(0)
const isDragging = ref(false)
const startX = ref(0)
const currentX = ref(0)
const currentIndex = ref(0)
const totalItems = ref(0)
const containerWidth = ref(0)
let autoPlayTimer = null

const maxIndex = computed(() => {
  return Math.max(0, totalItems.value - props.itemsPerView)
})

function handleTouchStart(e) {
  isDragging.value = true
  startX.value = e.touches ? e.touches[0].clientX : e.clientX
  currentX.value = startX.value
  
  if (props.autoPlay) {
    stopAutoPlay()
  }
}

function handleTouchMove(e) {
  if (!isDragging.value) return
  
  const x = e.touches ? e.touches[0].clientX : e.clientX
  currentX.value = x
  
  const diff = currentX.value - startX.value
  translateX.value = diff
}

function handleTouchEnd() {
  if (!isDragging.value) return
  
  isDragging.value = false
  
  const diff = currentX.value - startX.value
  const threshold = 50 // العتبة للتنقل
  
  if (diff > threshold) {
    // السحب لليمين - العودة للسابق
    prevSlide()
  } else if (diff < -threshold) {
    // السحب لليسار - الذهاب للتالي
    nextSlide()
  } else {
    // العودة للوضع الحالي
    resetPosition()
  }
  
  if (props.autoPlay) {
    startAutoPlay()
  }
}

function nextSlide() {
  if (currentIndex.value < maxIndex.value) {
    currentIndex.value++
    emit('swipe-left', currentIndex.value)
  } else {
    // العودة للبداية
    currentIndex.value = 0
  }
  resetPosition()
  emit('slide-change', currentIndex.value)
}

function prevSlide() {
  if (currentIndex.value > 0) {
    currentIndex.value--
    emit('swipe-right', currentIndex.value)
  } else {
    // الذهاب للنهاية
    currentIndex.value = maxIndex.value
  }
  resetPosition()
  emit('slide-change', currentIndex.value)
}

function goToSlide(index) {
  currentIndex.value = index
  resetPosition()
  emit('slide-change', currentIndex.value)
}

function resetPosition() {
  const itemWidth = containerWidth.value / props.itemsPerView
  translateX.value = -currentIndex.value * itemWidth
}

function updateTotalItems() {
  const content = document.querySelector('.swipe-content')
  if (content) {
    totalItems.value = content.children.length
    containerWidth.value = content.offsetWidth
    resetPosition()
  }
}

function startAutoPlay() {
  if (autoPlayTimer) return
  autoPlayTimer = setInterval(() => {
    nextSlide()
  }, props.autoPlayInterval)
}

function stopAutoPlay() {
  if (autoPlayTimer) {
    clearInterval(autoPlayTimer)
    autoPlayTimer = null
  }
}

onMounted(() => {
  updateTotalItems()
  window.addEventListener('resize', updateTotalItems)
  
  if (props.autoPlay) {
    startAutoPlay()
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', updateTotalItems)
  stopAutoPlay()
})

// توفير الوصول للوظائف للمكونات الأب
defineExpose({
  nextSlide,
  prevSlide,
  goToSlide,
  currentIndex
})
</script>

<style scoped>
.swipe-nav {
  position: relative;
  overflow: hidden;
  width: 100%;
}

.swipe-content {
  display: flex;
  width: 100%;
  cursor: grab;
  user-select: none;
  -webkit-user-select: none;
  touch-action: pan-y;
}

.swipe-content:active {
  cursor: grabbing;
}

.swipe-indicators {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-top: 16px;
  padding: 8px;
}

.indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #E2E8F0;
  cursor: pointer;
  transition: all 0.3s ease;
}

.indicator.active {
  background: #1E3A5F;
  transform: scale(1.2);
}

.indicator:hover {
  background: #4A6FA5;
}

@media (max-width: 480px) {
  .indicator {
    width: 6px;
    height: 6px;
  }
  
  .swipe-indicators {
    gap: 6px;
    margin-top: 12px;
  }
}
</style>
