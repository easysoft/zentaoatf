<template>
  <div class="result-list">
    <List compact divider>
      <ListItem
        icon="checkmark-circle"
        iconClass="text-green"
        title="current_test.php"
        trailingText="13:00"
        trailingTextClass="muted small"
      />
      <ListItem
        icon="close-circle"
        iconClass="text-red"
        title="1_string_match.php"
        trailingText="12:40"
        trailingTextClass="muted small"
      />
      <ListItem
        icon="close-circle"
        iconClass="text-red"
        title="chrome.php"
        trailingText="12:10"
        trailingTextClass="muted small"
      />
      <ListItem
        icon="checkmark-circle"
        iconClass="text-green"
        title="node_test.js"
        trailingText="11:30"
        trailingTextClass="muted small"
      />
      <ListItem
        icon="checkmark-circle"
        iconClass="text-green"
        title="批量：5个文件"
        trailingText="09:30"
        trailingTextClass="muted small"
      />
      <ListItem
        icon="close-circle"
        iconClass="text-red"
        title="批量：20个文件"
        trailingText="08:30"
        trailingTextClass="muted small"
      />
      <ListItem
        icon="checkmark-circle"
        iconClass="text-green"
        title="your_big_test.bat"
        trailingText="昨天 08:30"
        trailingTextClass="muted small"
      />
    </List>
  </div>
</template>

<script setup lang="ts">
import List from './List.vue';
import ListItem from './ListItem.vue';
</script>
