package main

import (
	"fmt"

	"demo/redis-demo/2.02.tag/tag"
	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	tg := tag.NewTag(client)

	// 添加标签
	tg.Add("EleutherAI_gpt-neo-1.3B", []string{"特征提取", "英语", "100B<n<1T"})
	tg.Add("stabilityai_stable-diffusion-2-1", []string{"文本生成图片", "英语", "10B<n<100B"})

	// 获取标签
	tags, _ := tg.GetTagsByTarget("EleutherAI_gpt-neo-1.3B")
	fmt.Println("Tags for EleutherAI_gpt-neo-1.3B:", tags)

	// 根据标签获取目标
	targets, _ := tg.GetTargetByTags([]string{"英语"})
	fmt.Println("Targets with 英语:", targets)

	// 缓存版获取目标
	cachedTargets, _ := tg.GetCachedTargetByTags([]string{"英语"})
	fmt.Println("Cached targets:", cachedTargets)
}
