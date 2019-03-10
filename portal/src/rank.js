import moment from 'moment';

function idToChar(id) {
  return "ABCDEFGHIJKLMNOPQRSTUVWXYZ".charAt(id);
}

export const Homework = {
  calc: function (user, problem, submission, contest) {
    moment.locale("zh-CN");
    var records = [];
    var records_rev = {};
    var probs = [];
    var probs_rev = {};
    var probs_rrev = {};
    var title = [{ text: "Rank", field: "rank" }, { text: "Solve", field: "solve" }];
    let realId = 0;
    for (let i = 0; i < problem.length; i++) {
      probs_rev[problem[i].uid] = realId;
      probs_rrev[realId] = problem[i].uid;
      realId++;
    }
    for (let i = 0; i < user.length; i++) {
      records_rev[user[i].uid] = i;
      if (user[i].name.length > 16) {
        user[i].name = user[i].name.substr(0, 14) + '...';
      }
      var temp = [];
      for (let j = 0; j < realId; j++) {
        temp.push({ uid: probs_rrev[j], att: 0, export: '未提交', time: idToChar(j), ac: 0, state: 'NOT_ATT', append: '' });
      }
      records.push({
        user: user[i],
        solve: 0,
        penalty: 0,
        last: 0,
        aftersolve: 0,
        result: temp
      });
    }
    for (let i = 0; i < problem; i++) {
      probs.push({
        att: 0,
        ac: 0,
        firstblood: -1
      });
    }
    submission.sort((a, b) => (a.submit_time < b.submit_time ? -1 : 1));
    for (var i = 0; i < submission.length; i++) {
      if (records_rev[submission[i].user_uid] == undefined) continue;
      var team = records_rev[submission[i].user_uid];
      var prob = probs_rev[submission[i].problem_uid];
      var verd = submission[i].verdict;
      
      if (prob === undefined) continue;
      if (records[team].result[prob].ac) continue;
      if (verd == 'AC') {
        if ((new Date(submission[i].submit_time)).valueOf() > contest.end_time.valueOf()) {
          records[team].result[prob].state = 'AC_AFTER';
          records[team].result[prob].time = moment(submission[i].submit_time).fromNow() + " (overdue)";
          records[team].result[prob].export = "通过（overdue）";
          records[team].aftersolve += 1;
        } else {
          records[team].result[prob].state = 'AC';
          records[team].solve += 1;
          records[team].result[prob].export = "通过";
          records[team].result[prob].time = moment(submission[i].submit_time).fromNow();
        }
        records[team].last = (new Date(submission[i].submit_time)).valueOf();
        records[team].result[prob].ac = 1;
      } else if (verd == 'PENDING') {
        records[team].result[prob].state = 'PENDING';
        records[team].result[prob].att += 1;
        records[team].result[prob].export = "未通过";
        records[team].result[prob].time = (records[team].result[prob].att + " 次尝试");
      } else {
        records[team].result[prob].state = 'NO';
        records[team].result[prob].append = '';
        records[team].result[prob].att += 1;
        records[team].result[prob].export = "未通过";
        records[team].result[prob].time = (records[team].result[prob].att + " 次尝试");
      }
    }
    //------//
    var sorted = [];
    for (var i = 0; i < records.length; i++) { sorted.push(records[i]); }
    function cmp(x, y) {
      if (x.solve != y.solve) return y.solve - x.solve;
      else if (x.aftersolve != y.aftersolve) return y.aftersolve - x.aftersolve;
      else if (x.last != y.last) return x.last - y.last;
      return x.user.name < y.user.name ? -1 : 1;
    }
    function cmp2(x, y) {
      if (x.solve != y.solve) return y.solve - x.solve;
      return 0;
    }
    sorted.sort(cmp);
    let realRank = 0, last_valid = -1;
    for (let i = 0; i < sorted.length; i++) {
      realRank++;
      if (last_valid == -1 || cmp2(sorted[last_valid], sorted[i]) < 0) {
        sorted[i].rank = realRank;
      } else {
        sorted[i].rank = sorted[last_valid].rank;
      }
      last_valid = i;
    }
    return { title: title, rank: sorted, problem: probs };
  }
}

export const ICPC = {
  calc: function (user, problem, submission, contest) {
    moment.locale("zh-CN");
    var records = [];
    var records_rev = {};
    var probs = [];
    var probs_rev = {};
    var probs_rrev = {};
    var title = [{ text: "Rank", field: "rank" }, { text: "Solve", field: "solve" }, { text: "Penalty", field: "penalty" }];
    contest = JSON.parse(JSON.stringify(contest));
    contest.start_time = new Date(contest.start_time);
    contest.end_time = new Date(contest.end_time);
    contest.freeze_time = new Date(contest.freeze_time);
    
    let realId = 0;
    for (let i = 0; i < problem.length; i++) {
      probs_rev[problem[i].uid] = realId;
      probs_rrev[realId] = problem[i].uid;
      realId++;
    }
    for (let i = 0; i < user.length; i++) {
      records_rev[user[i].uid] = i;
      if (user[i].name.length > 16) {
        user[i].name = user[i].name.substr(0, 14) + '...';
      }
      var temp = [];
      for (let j = 0; j < realId; j++) {
        temp.push({ uid: probs_rrev[j], att: 0, pen: 0, time: idToChar(j), ac: 0, state: 'NOT_ATT', append: '' });
      }
      records.push({
        user: user[i],
        solve: 0,
        penalty: 0,
        last: 0,
        result: temp
      });
    }
    for (let i = 0; i < realId; i++) {
      probs.push({
        att: 0,
        ac: 0,
        firstblood: -1
      });
    }
    submission.sort((a, b) => (a.submit_time < b.submit_time ? -1 : 1));
    for (var i = 0; i < submission.length; i++) {
      if (records_rev[submission[i].user_uid] == undefined) continue;
      var team = records_rev[submission[i].user_uid];
      var prob = probs_rev[submission[i].problem_uid];
      var verd = submission[i].verdict;
      
      if (prob === undefined) continue;
      if (records[team].result[prob].ac) continue;
      if (verd == 'AC') {
        if ((new Date(submission[i].submit_time)).valueOf() > contest.end_time.valueOf()) {
          records[team].result[prob].state = 'AC_AFTER';
          //records[team].result[prob].time = moment(submission[i].submit_time).fromNow();
        } else {
          records[team].result[prob].state = 'AC';
          records[team].penalty += records[team].result[prob].att * 20 + Math.floor(((new Date(submission[i].submit_time)).valueOf() - contest.start_time.valueOf()) / 1000 / 60)
          records[team].solve += 1;
          records[team].result[prob].time = (records[team].result[prob].att + records[team].result[prob].pen + 1) + " - " + Math.floor(((new Date(submission[i].submit_time)).valueOf() - contest.start_time.valueOf()) / 1000 / 60);
          probs[prob].ac++;
          probs[prob].att++;
          if (probs[prob].firstblood == -1) {
            probs[prob].firstblood = Math.floor(((new Date(submission[i].submit_time)).valueOf() - contest.start_time.valueOf()) / 1000 / 60);
          }
          if (Math.floor(((new Date(submission[i].submit_time)).valueOf() - contest.start_time.valueOf()) / 1000 / 60) == probs[prob].firstblood) {
            records[team].result[prob].state = 'AC_FIRST';
          }
        }
        records[team].result[prob].ac = 1;
      } else if (verd == 'PENDING') {
        if ((new Date(submission[i].submit_time)).valueOf() > contest.end_time.valueOf()) {
          continue;
        }
        probs[prob].att++;
        records[team].result[prob].state = 'PENDING';
        records[team].result[prob].pen += 1;
        records[team].result[prob].time = (records[team].result[prob].att + records[team].result[prob].pen) + " - " + Math.floor(((new Date(submission[i].submit_time)).valueOf() - contest.start_time.valueOf()) / 1000 / 60);
      } else {
        if ((new Date(submission[i].submit_time)).valueOf() > contest.end_time.valueOf()) {
          continue;
        }
        probs[prob].att++;
        records[team].result[prob].state = 'NO';
        records[team].result[prob].append = '';
        records[team].result[prob].att += 1;
        records[team].result[prob].time = (records[team].result[prob].att + records[team].result[prob].pen) + " - " + Math.floor(((new Date(submission[i].submit_time)).valueOf() - contest.start_time.valueOf()) / 1000 / 60);
        //records[team].result[prob].time = (records[team].result[prob].att + " 次尝试");
      }
    }
    //------//
    var sorted = [];
    for (var i = 0; i < records.length; i++) { sorted.push(records[i]); }
    function cmp(x, y) {
      if (x.solve != y.solve) return y.solve - x.solve;
      else if (x.penalty != y.penalty) return x.penalty - y.penalty;
      else if (x.last != y.last) return x.last - y.last;
      return x.user.name < y.user.name ? -1 : 1;
    }
    function cmp2(x, y) {
      if (x.solve != y.solve) return y.solve - x.solve;
      else if (x.penalty != y.penalty) return x.penalty - y.penalty;
      else if (x.last != y.last) return x.last - y.last;
      return 0;
    }
    sorted.sort(cmp);
    let realRank = 0, last_valid = -1;
    for (let i = 0; i < sorted.length; i++) {
      realRank++;
      if (last_valid == -1 || cmp2(sorted[last_valid], sorted[i]) < 0) {
        sorted[i].rank = realRank;
      } else {
        sorted[i].rank = sorted[last_valid].rank;
      }
      last_valid = i;
    }
    return { title: title, rank: sorted, problem: probs };
  }
}

// 更适合个人比赛的ICPC模式
// 每题通过时间取Max而非求和
// 罚时为6分钟（约等于20分钟/3人）而非20分钟
export const ICPCInd = {
  calc: function (user, problem, submission, contest) {
    moment.locale("zh-CN");
    var records = [];
    var records_rev = {};
    var probs = [];
    var probs_rev = {};
    var probs_rrev = {};
    var title = [{ text: "Rank", field: "rank" }, { text: "Solve", field: "solve" }, { text: "Penalty", field: "penalty" }];
    let realId = 0;
    for (let i = 0; i < problem.length; i++) {
      probs_rev[problem[i].uid] = realId;
      probs_rrev[realId] = problem[i].uid;
      realId++;
    }
    for (let i = 0; i < user.length; i++) {
      records_rev[user[i].uid] = i;
      if (user[i].name.length > 16) {
        user[i].name = user[i].name.substr(0, 14) + '...';
      }
      var temp = [];
      for (let j = 0; j < realId; j++) {
        temp.push({ uid: probs_rrev[j], att: 0, pen: 0, time: idToChar(j), ac: 0, state: 'NOT_ATT', append: '' });
      }
      records.push({
        user: user[i],
        solve: 0,
        penalty: 0,
        last: 0,
        result: temp
      });
    }
    for (let i = 0; i < realId; i++) {
      probs.push({
        att: 0,
        ac: 0,
        firstblood: -1
      });
    }
    submission.sort((a, b) => (a.submit_time < b.submit_time ? -1 : 1));
    for (var i = 0; i < submission.length; i++) {
      if (records_rev[submission[i].user_uid] == undefined) continue;
      var team = records_rev[submission[i].user_uid];
      var prob = probs_rev[submission[i].problem_uid];
      var verd = submission[i].verdict;
      
      if (prob === undefined) continue;
      if (records[team].result[prob].ac) continue;
      if (verd == 'AC') {
        if ((new Date(submission[i].submit_time)).valueOf() > contest.end_time.valueOf()) {
          records[team].result[prob].state = 'AC_AFTER';
          //records[team].result[prob].time = moment(submission[i].submit_time).fromNow();
        } else {
          records[team].result[prob].state = 'AC';
          records[team].penalty += records[team].result[prob].att * 6;
          records[team].last = Math.floor(((new Date(submission[i].submit_time)).valueOf() - contest.start_time.valueOf()) / 1000 / 60);
          records[team].solve += 1;
          records[team].result[prob].time = (records[team].result[prob].att + records[team].result[prob].pen + 1) + " - " + Math.floor(((new Date(submission[i].submit_time)).valueOf() - contest.start_time.valueOf()) / 1000 / 60);
          if (probs[prob].firstblood == -1) {
            probs[prob].firstblood = Math.floor(((new Date(submission[i].submit_time)).valueOf() - contest.start_time.valueOf()) / 1000 / 60);
          }
          if (Math.floor(((new Date(submission[i].submit_time)).valueOf() - contest.start_time.valueOf()) / 1000 / 60) == probs[prob].firstblood) {
            records[team].result[prob].state = 'AC_FIRST';
          }
        }
        records[team].result[prob].ac = 1;
      } else if (verd == 'PENDING') {
        if ((new Date(submission[i].submit_time)).valueOf() > contest.end_time.valueOf()) {
          continue;
        }
        records[team].result[prob].state = 'PENDING';
        records[team].result[prob].pen += 1;
        records[team].result[prob].time = (records[team].result[prob].att + records[team].result[prob].pen) + " - " + Math.floor(((new Date(submission[i].submit_time)).valueOf() - contest.start_time.valueOf()) / 1000 / 60);
      } else {
        if ((new Date(submission[i].submit_time)).valueOf() > contest.end_time.valueOf()) {
          continue;
        }
        records[team].result[prob].state = 'NO';
        records[team].result[prob].append = '';
        records[team].result[prob].att += 1;
        records[team].result[prob].time = (records[team].result[prob].att + records[team].result[prob].pen) + " - " + Math.floor(((new Date(submission[i].submit_time)).valueOf() - contest.start_time.valueOf()) / 1000 / 60);
        //records[team].result[prob].time = (records[team].result[prob].att + " 次尝试");
      }
    }
    //------//
    var sorted = [];
    for (var i = 0; i < records.length; i++) {
      records[i].penalty += records[i].last;
      sorted.push(records[i]);
    }
    function cmp(x, y) {
      if (x.solve != y.solve) return y.solve - x.solve;
      else if (x.penalty != y.penalty) return x.penalty - y.penalty;
      else if (x.last != y.last) return x.last - y.last;
      return x.user.name < y.user.name ? -1 : 1;
    }
    function cmp2(x, y) {
      if (x.solve != y.solve) return y.solve - x.solve;
      else if (x.penalty != y.penalty) return x.penalty - y.penalty;
      else if (x.last != y.last) return x.last - y.last;
      return 0;
    }
    sorted.sort(cmp);
    let realRank = 0, last_valid = -1;
    for (let i = 0; i < sorted.length; i++) {
      realRank++;
      if (last_valid == -1 || cmp2(sorted[last_valid], sorted[i]) < 0) {
        sorted[i].rank = realRank;
      } else {
        sorted[i].rank = sorted[last_valid].rank;
      }
      last_valid = i;
    }
    return { title: title, rank: sorted, problem: probs };
  }
}