"""
12-Agent Game with Multi-Phase Mechanics

PHASES:
1. Fog (Rounds 1-4): Votes are secret, only damage shown
2. Reveal (Rounds 5-9): Full vote logs shown
3. Endgame (Last 4 agents): 4 damage + one-time ultimates
"""

import random
from collections import defaultdict

STRATEGIES = {
    'entropy_shift': ('Alpha', '混沌不可预测'),
    'grudge_keeper': ('Beta', '记仇报复'),
    'poison_leader': ('Gamma', '擒贼擒王'),
    'survivor_hunter': ('Delta', '收割残血'),
    'bandwagon_breaker': ('Epsilon', '反主流'),
    'chaos_agent': ('Zeta', '纯混沌'),
    'kingmaker': ('Eta', '造王者'),
    'mirror_twin': ('Theta', '镜像防御'),
    'assassin': ('Iota', '蓄力刺客'),
    'prophecy': ('Kappa', '预测大师'),
}

class Agent:
    def __init__(self, name, strategy, description, max_hp=100):
        self.name = name
        self.strategy = strategy
        self.description = description
        self.hp = max_hp
        self.votes_received = 0
        self.is_alive = True
        self.assassin_charge = 0
        self.last_target = None
        self.ultimates_used = False
    
    def select_vote(self, agents, history, phase, round_num):
        if not self.is_alive:
            return None
        
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
        
        vote_fn = strategy_map.get(self.strategy, self.random_vote)
        return vote_fn(agents, history, phase)
    
    def get_living_agents(self, agents):
        return [a for a in agents if a.is_alive and a.name != self.name]
    
    def use_ultimate(self, agents, phase):
        """终极技能 - 终局阶段使用"""
        if self.ultimates_used or phase != 'endgame':
            return None
        
        living = self.get_living_agents(agents)
        if not living:
            return None
        
        # 根据策略选择终极技能目标
        if self.strategy == 'poison_leader':
            # 双重伤害
            living.sort(key=lambda a: -a.hp)
            self.ultimates_used = True
            return {'type': 'double_damage', 'target': living[0].name}
        
        elif self.strategy == 'survivor_hunter':
            # 终结最低血量
            living.sort(key=lambda a: a.hp)
            self.ultimates_used = True
            return {'type': 'execute', 'target': living[0].name}
        
        elif self.strategy == 'grudge_keeper':
            # 对准仇人双重伤害
            offender = None
            for v in reversed(history):
                if v['target'] == self.name and v['voter'] != self.name:
                    if any(a.is_alive and a.name == v['voter'] for a in agents):
                        offender = v['voter']
                        break
            if offender:
                self.ultimates_used = True
                return {'type': 'revenge', 'target': offender}
        
        # 默认：对自己使用护盾
        self.hp += 10
        self.ultimates_used = True
        return {'type': 'heal', 'amount': 10}
    
    def get_quote(self, round_num, phase):
        phase_names = {'fog': '迷雾', 'reveal': '揭露', 'endgame': '终局'}
        quotes = {
            'fog': [
                "迷雾笼罩，没有人知道真相。",
                "暗中观察，等待时机。",
                "投票隐藏，行动致命。",
            ],
            'reveal': [
                "真相大白，是时候算账了。",
                "所有投票都已记录。",
                "政治游戏现在开始。",
            ],
            'endgame': [
                "这是最后一战！",
                "终极技能，扭转乾坤。",
                "不是生存，就是死亡。",
            ],
        }
        q = quotes.get(phase, quotes['fog'])
        return q[round_num % len(q)]
    
    # 策略实现
    def entropy_shift_vote(self, agents, history, phase):
        if random.random() < 0.4:
            return self.random_vote(agents, history)
        return self.analyze_voting_patterns(agents, history)
    
    def grudge_keeper_vote(self, agents, history, phase):
        offender_count = defaultdict(int)
        for v in history:
            if v['target'] == self.name and v['voter'] != self.name:
                offender_count[v['voter']] += 1
        if offender_count:
            target = max(offender_count, key=offender_count.get)
            if any(a.is_alive and a.name == target for a in agents):
                return target
        return self.random_vote(agents, history)
    
    def poison_leader_vote(self, agents, history, phase):
        living = self.get_living_agents(agents)
        if not living:
            return ""
        living.sort(key=lambda a: (-a.hp, a.name))
        return living[0].name
    
    def survivor_hunter_vote(self, agents, history, phase):
        living = self.get_living_agents(agents)
        if not living:
            return ""
        living.sort(key=lambda a: a.hp)
        return living[0].name
    
    def bandwagon_breaker_vote(self, agents, history, phase):
        vote_count = defaultdict(int)
        for v in history:
            vote_count[v['target']] += 1
        if vote_count:
            popular = max(vote_count, key=vote_count.get)
            if vote_count[popular] >= 1:
                living = self.get_living_agents(agents)
                for a in living:
                    if a.name != popular:
                        return a.name
        return self.random_vote(agents, history)
    
    def chaos_agent_vote(self, agents, history, phase):
        return self.random_vote(agents, history)
    
    def kingmaker_vote(self, agents, history, phase):
        living = self.get_living_agents(agents)
        if not living:
            return ""
        living.sort(key=lambda a: (-a.hp, a.name))
        if len(living) >= 2 and living[0].name != self.name:
            return living[1].name
        return living[0].name if living else ""
    
    def mirror_twin_vote(self, agents, history, phase):
        for v in reversed(history):
            if v['target'] == self.name and v['voter'] != self.name:
                if any(a.is_alive and a.name == v['voter'] for a in agents):
                    return v['voter']
        return self.random_vote(agents, history)
    
    def assassin_vote(self, agents, history, phase):
        CHARGE_NEEDED = 2
        if self.assassin_charge < CHARGE_NEEDED:
            self.assassin_charge += 1
            return ""
        self.assassin_charge = 0
        living = self.get_living_agents(agents)
        if not living:
            return ""
        living.sort(key=lambda a: a.hp)
        return living[0].name
    
    def prophecy_vote(self, agents, history, phase):
        if not history:
            return self.random_vote(agents, history)
        last_round = max(v['round'] for v in history)
        counts = defaultdict(int)
        for v in history:
            if v['round'] == last_round:
                counts[v['target']] += 1
        if counts:
            predicted = max(counts, key=counts.get)
            if predicted != self.name and any(a.is_alive and a.name == predicted for a in agents):
                return predicted
        return self.random_vote(agents, history)
    
    def analyze_voting_patterns(self, agents, history):
        return self.random_vote(agents, history)
    
    def random_vote(self, agents, history):
        living = self.get_living_agents(agents)
        if not living:
            return ""
        return random.choice(living).name
    
    def take_damage(self, damage):
        self.hp -= damage
        self.votes_received += damage
        if self.hp <= 0:
            self.hp = 0
            self.is_alive = False


class GameEngine:
    def __init__(self, max_rounds=15, base_hp=100, attack_damage=3):
        self.max_rounds = max_rounds
        self.base_hp = base_hp
        self.attack_damage = attack_damage
    
    def run_game(self):
        agents = []
        for strat, (name, desc) in STRATEGIES.items():
            agents.append(Agent(name, strat, desc, self.base_hp))
        
        history = []
        scenes = []
        
        for round_num in range(1, self.max_rounds + 1):
            alive = [a for a in agents if a.is_alive]
            
            if len(alive) <= 1:
                break
            
            # 确定阶段
            if len(alive) > 8:
                phase = 'fog'
                damage = self.attack_damage
            elif len(alive) > 4:
                phase = 'reveal'
                damage = self.attack_damage
            else:
                phase = 'endgame'
                damage = self.attack_damage + 1  # 终局伤害+1
            
            phase_names = {'fog': '🌫️ 迷雾', 'reveal': '👁️ 揭露', 'endgame': '🔥 终局'}
            
            scenes.append(f"\n{'='*70}")
            scenes.append(f"🎮 第 {round_num} 轮 - {phase_names[phase]} - {len(alive)} 人存活")
            scenes.append(f"{'='*70}")
            
            # 存活状态
            alive.sort(key=lambda a: (-a.hp, a.name))
            status = " | ".join([f"{a.name}({a.hp})" for a in alive])
            scenes.append(f"📊 状态: {status}")
            
            # 终极技能使用情况
            ultimates_this_round = []
            for agent in alive:
                if phase == 'endgame' and not agent.ultimates_used:
                    ult = agent.use_ultimate(agents, phase)
                    if ult:
                        ultimates_this_round.append((agent.name, ult))
            
            if ultimates_this_round:
                for name, ult in ultimates_this_round:
                    if ult['type'] == 'double_damage':
                        scenes.append(f"⚡ {name} 释放终极技能: 对 {ult['target']} 双倍伤害！")
                    elif ult['type'] == 'execute':
                        scenes.append(f"⚡ {name} 释放终极技能: 终结 {ult['target']}！")
                    elif ult['type'] == 'revenge':
                        scenes.append(f"⚡ {name} 释放终极技能: 复仇 {ult['target']}！")
                    elif ult['type'] == 'heal':
                        scenes.append(f"💚 {name} 释放终极技能: 恢复 {ult['amount']} HP！")
            
            # 投票阶段
            votes = {}
            dialogues = []
            for agent in alive:
                target = agent.select_vote(agents, history, phase, round_num)
                if target:
                    votes[target] = votes.get(target, 0) + 1
                    history.append({'round': round_num, 'voter': agent.name, 'target': target})
                    quote = agent.get_quote(round_num, phase)
                    dialogues.append(f"  [{name}→{target}]: \"{quote}\"")
            
            scenes.append(f"\n💬 {phase_names[phase]}阶段对话:")
            scenes.extend(dialogues)
            
            # 投票结果 - 迷雾阶段隐藏投票者
            if phase == 'fog':
                scenes.append(f"\n🗳️ 投票结果: {votes} (攻击者隐藏)")
            else:
                scenes.append(f"\n🗳️ 投票结果: {votes}")
            
            # 攻击结算
            deaths = []
            for target, count in votes.items():
                dmg = damage * count
                for agent in agents:
                    if agent.name == target and agent.is_alive:
                        agent.take_damage(dmg)
                        scenes.append(f"  ⚔️ {target} 受到 {dmg} 伤害 ({count}票 × {damage})")
                        if not agent.is_alive:
                            deaths.append(target)
            
            if deaths:
                scenes.append(f"\n💀 淘汰: {', '.join(deaths)}")
            
            # 检查游戏结束
            alive_now = [a for a in agents if a.is_alive]
            if len(alive_now) <= 1:
                break
        
        survivors = [a for a in agents if a.is_alive]
        return {
            'winner': survivors[0].name if survivors else None,
            'rounds': round_num,
            'scenes': scenes,
            'final_hp': {a.name: a.hp for a in agents},
            'death_order': [a.name for a in agents if not a.is_alive],
            'phase_stats': {
                'fog_rounds': sum(1 for r in range(1, round_num+1) if sum(1 for a in agents if a.is_alive and sum(1 for _ in []) > 8)),
            }
        }


# ============ 运行游戏 ============
print("="*70)
print("🎮 12 人三阶段游戏 - 迷雾 → 揭露 → 终局")
print("="*70)
print("\n阶段规则:")
print("  🌫️ 迷雾 (8+人): 投票隐藏，只显示伤害")
print("  👁️ 揭露 (4-8人): 显示完整投票记录")
print("  🔥 终局 (≤4人): 伤害+1，每人一次终极技能")

print("\n" + "="*70)
print("🏆 游戏开始!")
print("="*70)

engine = GameEngine(max_rounds=15, base_hp=100, attack_damage=3)
result = engine.run_game()

# 打印所有场景
for scene in result['scenes']:
    print(scene)

# 结果
print("\n" + "="*70)
print("🎉 游戏结束!")
print("="*70)
print(f"\n🏆 冠军: {result['winner']}")
print(f"📅 总回合数: {result['rounds']}")
print(f"💀 淘汰顺序: {' → '.join(result['death_order'])}")

# 统计
print("\n📊 最终状态:")
for name, hp in sorted(result['final_hp'].items(), key=lambda x: -x[1]):
    status = "🏆 冠军" if hp > 0 and name == result['winner'] else ("💀 淘汰" if hp == 0 else "存活")
    print(f"  {name:8}: {hp:3} HP {status}")
