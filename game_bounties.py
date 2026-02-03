"""
12-Agent Game with Bounty System

MECHANICS:
- Each round: 2 public bounties appear
- Completing bounty = rewards (HP, extra vote, etc.)
- First to X influence points wins (or last standing)
"""

import random
from collections import defaultdict

STRATEGIES = {
    'entropy_shift': ('Alpha', '混沌'),
    'grudge_keeper': ('Beta', '记仇'),
    'poison_leader': ('Gamma', '擒王'),
    'survivor_hunter': ('Delta', '收割'),
    'bandwagon_breaker': ('Epsilon', '反主流'),
    'chaos_agent': ('Zeta', '混沌'),
    'kingmaker': ('Eta', '造王'),
    'mirror_twin': ('Theta', '镜像'),
    'assassin': ('Iota', '刺客'),
    'prophecy': ('Kappa', '预测'),
}

class Agent:
    def __init__(self, name, strategy, max_hp=100):
        self.name = name
        self.strategy = strategy
        self.hp = max_hp
        self.is_alive = True
        self.assassin_charge = 0
        self.influence_points = 0  # 新：影响力积分
        self.bounty_completed = 0   # 新：完成的赏金
    
    def select_vote(self, agents, history, bounties):
        if not self.is_alive:
            return None
        living = [a for a in agents if a.is_alive and a.name != self.name]
        if not living:
            return ""
        
        # 策略实现
        strategy_map = {
            'entropy_shift': self.entropy_shift_vote,
            'grudge_keeper': self.grudge_keeper_vote,
            'poison_leader': self.poison_leader_vote,
            'survivor_hunter': self.survivor_hunter_vote,
            'bandwagon_breaker': self.bandwagon_breaker_vote,
            'chaos_agent': self.chaos_agent_vote,
            'kingmaker': self.kingmaker_vote,
            'mirror_twin': self.mirror_twin_vote,
            'assassin': self.assassin_vote,
            'prophecy': self.prophecy_vote,
        }
        return strategy_map.get(self.strategy, self.random_vote)(living, history, bounties)
    
    def entropy_shift_vote(self, living, history, bounties):
        # 赏金目标优先
        for bounty in bounties:
            if bounty['target'] in [a.name for a in living]:
                return bounty['target']
        return random.choice(living).name
    
    def grudge_keeper_vote(self, living, history, bounties):
        for v in reversed(history):
            if v['target'] == self.name and v['voter'] != self.name:
                if any(a.is_alive and a.name == v['voter'] for a in living):
                    return v['voter']
        for bounty in bounties:
            if bounty['target'] in [a.name for a in living]:
                return bounty['target']
        return random.choice(living).name
    
    def poison_leader_vote(self, living, history, bounties):
        living.sort(key=lambda a: -a.hp)
        return living[0].name
    
    def survivor_hunter_vote(self, living, history, bounties):
        living.sort(key=lambda a: a.hp)
        return living[0].name
    
    def bandwagon_breaker_vote(self, living, history, bounties):
        for bounty in bounties:
            if bounty['target'] in [a.name for a in living]:
                return bounty['target']
        vote_count = defaultdict(int)
        for v in history:
            vote_count[v['target']] += 1
        if vote_count:
            popular = max(vote_count, key=vote_count.get)
            for a in living:
                if a.name != popular:
                    return a.name
        return random.choice(living).name
    
    def chaos_agent_vote(self, living, history, bounties):
        for bounty in bounties:
            if bounty['target'] in [a.name for a in living] and random.random() < 0.3:
                return bounty['target']
        return random.choice(living).name
    
    def kingmaker_vote(self, living, history, bounties):
        living.sort(key=lambda a: -a.hp)
        if len(living) >= 2:
            return living[1].name
        return living[0].name if living else ""
    
    def mirror_twin_vote(self, living, history, bounties):
        for v in reversed(history):
            if v['target'] == self.name and v['voter'] != self.name:
                if any(a.is_alive and a.name == v['voter'] for a in living):
                    return v['voter']
        for bounty in bounties:
            if bounty['target'] in [a.name for a in living]:
                return bounty['target']
        return random.choice(living).name
    
    def assassin_vote(self, living, history, bounties):
        for bounty in bounties:
            if bounty['target'] in [a.name for a in living] and bounty['reward'] == 'extra_vote':
                return bounty['target']
        if self.assassin_charge < 2:
            self.assassin_charge += 1
            return ""
        self.assassin_charge = 0
        living.sort(key=lambda a: a.hp)
        return living[0].name
    
    def prophecy_vote(self, living, history, bounties):
        for bounty in bounties:
            if bounty['target'] in [a.name for a in living]:
                return bounty['target']
        if not history:
            return random.choice(living).name
        last_round = max(v['round'] for v in history)
        counts = defaultdict(int)
        for v in history:
            if v['round'] == last_round:
                counts[v['target']] += 1
        if counts:
            predicted = max(counts, key=counts.get)
            if predicted != self.name and any(a.is_alive and a.name == predicted for a in living):
                return predicted
        return random.choice(living).name
    
    def random_vote(self, living, history, bounties):
        return random.choice(living).name if living else ""
    
    def take_damage(self, damage):
        self.hp -= damage
        if self.hp <= 0:
            self.hp = 0
            self.is_alive = False


def generate_bounties(alive, round_num):
    """每轮生成 2 个赏金"""
    bounties = []
    targets = random.sample([a.name for a in alive if a.is_alive], min(2, len([a for a in alive if a.is_alive])))
    
    rewards = ['heal_5', 'extra_vote', 'pierce', 'influence_2']
    
    for target in targets:
        reward = random.choice(rewards)
        bounties.append({
            'target': target,
            'reward': reward,
            'description': {
                'heal_5': '+5 HP 完成赏金',
                'extra_vote': '下轮额外 1 票',
                'pierce': '伤害+1 (无视防御)',
                'influence_2': '+2 影响力'
            }[reward]
        })
    return bounties


def run_bounty_game():
    agents = [Agent(name, strat, 100) for strat, (name, _) in STRATEGIES.items()]
    history = []
    INFLUENCE_THRESHOLD = 15  # 影响力阈值
    
    for round_num in range(1, 20):
        alive = [a for a in agents if a.is_alive]
        if len(alive) <= 1:
            break
        
        # 检查影响力胜利
        for agent in alive:
            if agent.influence_points >= INFLUENCE_THRESHOLD:
                return {
                    'winner': f"{agent.name} (影响力胜利)",
                    'rounds': round_num,
                    'type': 'influence'
                }
        
        # 生成赏金
        bounties = generate_bounties(alive, round_num)
        
        # 投票
        votes = defaultdict(int)
        extra_votes = defaultdict(int)  # 额外票数
        
        for agent in alive:
            target = agent.select_vote(agents, history, bounties)
            if target:
                # 检查是否有额外票奖励
                has_extra = False
                for bounty in bounties:
                    if bounty['target'] == target and bounty['reward'] == 'extra_vote':
                        extra_votes[target] += 1
                        has_extra = True
                        break
                votes[target] += 1 + extra_votes[target]
                history.append({'round': round_num, 'voter': agent.name, 'target': target})
        
        # 结算赏金
        bounty_targets = [b['target'] for b in bounties]
        for target, count in votes.items():
            dmg = 3 * count
            for agent in agents:
                if agent.name == target and agent.is_alive:
                    agent.take_damage(dmg)
                    # 检查是否完成赏金
                    if target in bounty_targets:
                        for bounty in bounties:
                            if bounty['target'] == target:
                                # 发放奖励
                                if bounty['reward'] == 'heal_5':
                                    agent.hp = min(100, agent.hp + 5)
                                    agent.influence_points += 1
                                    agent.bounty_completed += 1
                                elif bounty['reward'] == 'extra_vote':
                                    agent.influence_points += 1
                                    agent.bounty_completed += 1
                                elif bounty['reward'] == 'pierce':
                                    agent.influence_points += 2
                                    agent.bounty_completed += 1
                                elif bounty['reward'] == 'influence_2':
                                    agent.influence_points += 2
                                    agent.bounty_completed += 1
    
    survivors = [a for a in agents if a.is_alive]
    if survivors:
        return {
            'winner': f"{survivors[0].name} (存活胜利)",
            'rounds': round_num,
            'type': 'survival',
            'final_hp': {a.name: a.hp for a in agents},
            'final_influence': {a.name: a.influence_points for a in agents},
        }
    return {'winner': '平局', 'rounds': round_num}


# 运行分析
print("="*70)
print("🎮 赏金系统游戏 - 100 次模拟")
print("="*70)
print("\n规则:")
print("  • 每轮 2 个公开赏金")
print("  • 完成赏金: +HP / 额外票 / +影响力")
print("  • 影响力先到 15 获胜")
print("  • 或最后存活者获胜")

wins = defaultdict(int)
win_types = {'influence': 0, 'survival': 0}
influence_stats = defaultdict(list)
hp_stats = defaultdict(list)

for i in range(100):
    result = run_bounty_game()
    winner = result['winner'].split(' ')[0]
    wins[winner] += 1
    win_types[result['type']] += 1
    
    for name, ip in result.get('final_influence', {}).items():
        influence_stats[name].append(ip)
    for name, hp in result.get('final_hp', {}).items():
        hp_stats[name].append(hp)

print(f"\n🏆 胜率排名:")
for name, count in sorted(wins.items(), key=lambda x: -x[1]):
    bar = "█" * (count // 2)
    print(f"  {name:8}: {count:3} 胜 ({count}%) {bar}")

print(f"\n📊 胜利类型:")
print(f"  影响力胜利: {win_types['influence']}")
print(f"  存活胜利: {win_types['survival']}")

print(f"\n📈 平均影响力 (越高=赏金猎手):")
for name, ips in sorted(influence_stats.items(), key=lambda x: sum(x[1])/len(x[1]), reverse=True):
    avg = sum(ips) / len(ips)
    bar = "▓" * int(avg)
    print(f"  {name:8}: {avg:.1f} {bar}")

print("\n" + "="*70)
print("💡 赏金系统洞察:")
print("  • 赏金改变了投票动机 (抢人头 > 随便打)")
print("  • 影响力系统提供了多种获胜路径")
print("  • 策略需要平衡: 抢赏金 vs 生存")
