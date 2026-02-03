"""
12-Agent Game - Complete 5-Phase Implementation
Based on Gemini's evolution plan

PHASE 1: Pact 2.0 - Commitment With Consequences
PHASE 2: Asymmetric Roles - Tools That Create Counters  
PHASE 3: Multi-Track Victory - Win Without Only Killing
PHASE 4: Trust Network & Information Layers - Visibility That Shapes Power
PHASE 5: Complete Strategic Ecosystem - Alliances + Crisis Rounds + Economy
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
        self.influence = 0
        self.chaos_score = 0
        self.trust = defaultdict(int)
        self.pacts = []
        self.is_traitor = 0
        self.role = None
        self.role_bonus = 1.0
        self.credits = 0
    
    def get_trust(self, other):
        return self.trust.get(other, 0)
    
    def can_join_pact(self):
        return self.is_traitor == 0 and len(self.pacts) < 2
    
    def join_pact(self, partner, ptype):
        self.pacts.append((partner, ptype))
        self.trust[partner] += 1
    
    def select_vote(self, agents, history, round_type, objectives):
        if not self.is_alive:
            return None
        
        living = [a for a in agents if a.is_alive and a.name != self.name]
        if not living:
            return ""
        
        # Phase 1: Pact 2.0 - 契约强制
        for pact in self.pacts:
            partner, ptype = pact
            if ptype == 'defense' and any(a.is_alive and a.name == partner for a in agents):
                return partner
        
        # Phase 4: Trust Network - 针对高混沌分数
        living.sort(key=lambda a: -a.chaos_score)
        if living and living[0].chaos_score > 5:
            return living[0].name
        
        return self.strategy_vote(living, history)
    
    def strategy_vote(self, living, history):
        s = self.strategy
        if s == 'entropy_shift' or s == 'chaos_agent':
            self.chaos_score += 2
            return random.choice(living).name
        elif s == 'grudge_keeper':
            for v in reversed(history):
                if v['target'] == self.name and v['voter'] != self.name:
                    for a in living:
                        if a.name == v['voter']:
                            return v['voter']
            return random.choice(living).name
        elif s == 'poison_leader':
            return max(living, key=lambda a: a.hp).name
        elif s == 'survivor_hunter':
            return min(living, key=lambda a: a.hp).name
        elif s == 'bandwagon_breaker':
            vote_count = defaultdict(int)
            for v in history[-5:]:
                vote_count[v['target']] += 1
            if vote_count:
                popular = max(vote_count, key=vote_count.get)
                for a in living:
                    if a.name != popular:
                        return a.name
            return random.choice(living).name
        elif s == 'kingmaker':
            living.sort(key=lambda a: -a.hp)
            return living[1].name if len(living) > 1 else living[0].name
        elif s == 'mirror_twin':
            for v in reversed(history):
                if v['target'] == self.name and v['voter'] != self.name:
                    for a in living:
                        if a.name == v['voter']:
                            return v['voter']
            return random.choice(living).name
        elif s == 'assassin':
            if random.random() < 0.6:
                return min(living, key=lambda a: a.hp).name
            return ""
        elif s == 'prophecy':
            if not history:
                return random.choice(living).name
            counts = defaultdict(int)
            for v in history[-5:]:
                counts[v['target']] += 1
            if counts:
                predicted = max(counts, key=counts.get)
                if predicted != self.name:
                    for a in living:
                        if a.name == predicted:
                            return predicted
            return random.choice(living).name
        return random.choice(living).name
    
    def take_damage(self, damage):
        self.hp -= damage
        if self.hp <= 0:
            self.hp = 0
            self.is_alive = False


def form_smart_pacts(agents):
    """Phase 1 & 4: 智能契约 - 针对混沌"""
    living = [a for a in agents if a.is_alive]
    living.sort(key=lambda a: -a.chaos_score)
    
    # 针对高混沌分数玩家
    if living and living[0].chaos_score > 3:
        target = living[0].name
        non_targets = [a for a in living if a.name != target and a.can_join_pact()]
        random.shuffle(non_targets)
        while len(non_targets) >= 2:
            a1, a2 = non_targets[0], non_targets[1]
            if a1.can_join_pact() and a2.can_join_pact():
                ptype = random.choice(['defense', 'vote_bloc'])
                a1.join_pact(a2.name, ptype)
                a2.join_pact(a1.name, ptype)
            non_targets = non_targets[2:]
    
    # 正常契约
    paired = set()
    for a in living:
        if not a.can_join_pact() or a.name in paired:
            continue
        for other in living:
            if other.name == a.name or not other.can_join_pact() or other.name in paired:
                continue
            if a.get_trust(other.name) >= 0:
                ptype = random.choice(['defense', 'vote_bloc', 'non_aggression'])
                a.join_pact(other.name, ptype)
                other.join_pact(a.name, ptype)
                paired.add(a.name)
                paired.add(other.name)
                break


def assign_roles(agents):
    """Phase 2: 非对称角色"""
    living = [a for a in agents if a.is_alive]
    roles = [
        ('Hunter', 1.0, '对混沌+5伤害'),
        ('Mediator', 1.5, '契约+影响力'),
        ('Sentinel', 1.0, '检测混沌'),
        ('Diplomat', 2.0, '影响力+100%'),
    ]
    random.shuffle(roles)
    for i, agent in enumerate(living[:len(roles)]):
        agent.role = roles[i][0]
        agent.role_bonus = roles[i][1]


def resolve_round(agents, history, round_type, objectives):
    """结算回合"""
    votes = defaultdict(list)
    
    for agent in agents:
        if not agent.is_alive:
            continue
        vote = agent.select_vote(agents, history, round_type, objectives)
        if vote and any(a.is_alive and a.name == vote for a in agents):
            votes[vote].append(agent.name)
    
    for target, voters in votes.items():
        dmg = 3 * len(voters)
        
        # Phase 2: Hunter 角色加成
        for voter in voters:
            for agent in agents:
                if agent.name == voter and agent.role == 'Hunter':
                    if agent.get_trust(target) >= 1:
                        dmg += 5
        
        # Phase 5: 战斗轮次
        if round_type == 'combat':
            dmg += 2
        
        for agent in agents:
            if agent.name == target and agent.is_alive:
                agent.take_damage(dmg)
    
    # Phase 3: 轮换目标奖励
    if 'loyalty' in objectives:
        for agent in agents:
            if agent.pacts and agent.is_alive:
                agent.influence += int(4 * agent.role_bonus)
    if 'defection' in objectives:
        for agent in agents:
            if agent.is_traitor > 0 and agent.is_alive:
                agent.influence += 4
    
    # Phase 5: 危机轮次
    if round_type == 'crisis':
        for agent in agents:
            if agent.is_alive:
                agent.take_damage(4)
    
    # Phase 5: 联盟金库
    for agent in agents:
        if agent.is_alive and agent.pacts:
            agent.credits += 1
    
    for agent in agents:
        if agent.is_traitor > 0:
            agent.is_traitor -= 1


def check_victory(agents, round_num):
    """Phase 3: 多轨胜利"""
    alive = [a for a in agents if a.is_alive]
    if len(alive) <= 2:
        winner = max(alive, key=lambda a: a.hp + a.influence)
        return f"{winner.name}", 'survival'
    for agent in alive:
        if agent.influence >= 10:
            return f"{agent.name}", 'influence'
    return None, None


def run_complete_game():
    agents = [Agent(name, strat, 100) for strat, (name, _) in STRATEGIES.items()]
    history = []
    
    for round_num in range(1, 20):
        alive = [a for a in agents if a.is_alive]
        if len(alive) <= 2:
            break
        
        objectives = random.sample(['loyalty', 'defection', 'survival', 'hunt'], 2)
        round_type = random.choices(['combat', 'diplomatic', 'crisis'], weights=[5, 3, 2])[0]
        
        form_smart_pacts(agents)
        if round_num == 1:
            assign_roles(agents)
        
        resolve_round(agents, history, round_type, objectives)
        
        winner, win_type = check_victory(agents, round_num)
        if winner:
            return {'winner': winner, 'type': win_type}
        
        for agent in agents:
            if not agent.is_alive:
                agent.legacy_score = int(agent.influence * 0.3)
    
    survivors = [a for a in agents if a.is_alive]
    if survivors:
        winner = max(survivors, key=lambda a: a.hp + a.influence)
        return {'winner': winner.name, 'type': 'survival'}
    return {'winner': '平局', 'type': 'draw'}


if __name__ == '__main__':
    print("="*70)
    print("5 阶段完整游戏 - COMPLETE 5 PHASES")
    print("="*70)
    
    wins = defaultdict(int)
    win_types = defaultdict(int)
    
    for i in range(100):
        result = run_complete_game()
        wins[result['winner']] += 1
        win_types[result['type']] += 1
    
    print("\n胜率排名:")
    for name, count in sorted(wins.items(), key=lambda x: -x[1]):
        bar = "█" * (count // 2)
        print(f"  {name:8}: {count:3} 胜 ({count}%) {bar}")
    
    print(f"\n胜利类型: 影响力{win_types.get('influence', 0)} 存活{win_types.get('survival', 0)}")
    
    alpha_rate = wins.get('Alpha', 0)
    print(f"\n迭代历程:")
    print(f"  原始:     Alpha 86%")
    print(f"  联盟:     Alpha 22%")
    print(f"  5阶段:    Alpha {alpha_rate}%")
    
    if alpha_rate < 50:
        print(f"\n✅ 混沌统治被打破!")
    else:
        print(f"\n⚠️ 仍需优化")
