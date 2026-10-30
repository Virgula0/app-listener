# Graph Report - app-listener  (2026-10-30)

## Corpus Check
- 562 files · ~1,788,665 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7309 nodes · 33396 edges · 252 communities (202 shown, 50 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2994 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7986a237`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- op
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- parseGuardEvent
- runWizard
- protectedPage
- batchInodeRows
- Network
- BinaryEntry
- IntegrationSuite
- guard_trust.bpf.c
- fyne.io/fyne/v2/widget.Label
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth_test.go
- FileTarget
- guard.bpf.c
- LibraryClosure
- go_pkg_os
- editView
- go_pkg_strings
- RefuseIdentityLines
- TestMonitorUseCaseLifecycle
- fileedit_model_test.go
- Guard
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- controlServer
- os.FileMode
- RefreshSection
- IntegrationSuite
- absPath
- monitorUnitTest
- cursor.go
- Candidate
- fakeDaemon
- stdio
- update_test.go
- fadeScroll
- backend.go
- Binary
- Client
- daemonconfig.go
- session
- .apply
- githubRelease
- mutedText
- GuardInodeKey
- Monitor
- Prereq
- .FilterExistingWhitelist
- inode_dev
- .filterEvent
- policyEditor
- unix_connect.c
- wm_test.go
- IntegrationSuite
- dialLiveSession
- runOfflineEdit
- wipe
- IntegrationSuite
- journal
- fcntl
- RawTarget
- LaunchRule
- codecheck/main.go
- NetworkMonitorRepository
- User
- process_vm_readv.c
- os.File
- recordFlow
- newTestTarget
- rcUser
- golang.org/x/sys/unix.Stat_t
- bpfstats/main.go
- exe_supersede.h
- loginServer
- vaultfs_test.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runView
- downloadAndVerify
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- hashFileShared
- root.go
- Well-known bypass classes
- hashCache
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- newTestVetter
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- highlight_test.go
- ConfinedEntry
- installPage
- Vault
- IntegrationSuite
- ListUsers
- deploy
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- UI
- Recorder
- WebSocket /ws client
- inet_connect.c
- install.sh
- mc_tag.h
- classify
- make test-integration (rootful Docker suite)
- IntegrationSuite
- catalogRefresher
- audit.go
- trace-app-libs.sh
- IntegrationSuite
- changelog.go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- planUpdaters
- serve.go
- liveSession
- github.com/charmbracelet/bubbles/textarea.Model
- Ledger
- appTheme
- GuardEvent
- buildTrustedSet
- os.FileInfo
- fakeStore
- bufio.Reader
- IntegrationSuite
- buildOneGuard
- dialog_test.go
- copytree.go
- confirmUpdate
- Installation
- backend_test.go
- graphify-refresh.sh
- RootPlacedName
- NetEventType
- secretJudge
- Row
- feed.go
- shell
- layout.go
- procstats.go
- ensurePortAnchors
- SelectOrphans
- Limitations
- RunWith
- IntegrationSuite
- TrustGuard
- net_tester/main.go
- sync.Once
- run
- NetGuardEvent
- fileedit_symlink_test.go
- How it works
- inodeChain
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- github.com/spf13/cobra.Command
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- IntegrationSuite
- Run
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- New
- pathCache
- sync.Mutex
- time.Time
- io.Writer
- EventType
- .populateWatchInodes
- monitor
- IntegrationSuite
- tempGrantSetup
- lockAndDeprovision
- .eventTypeAllowed
- KernelDev
- VouchedDevs
- IntegrationSuite
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- askSSHAgentUsers
- .requirePortDenied
- openLedger
- fakeVault
- downloadFile
- github.com/stretchr/testify/suite.Suite
- record
- ReadInstallCert
- supersede_probe.c
- github.com/charmbracelet/bubbletea.Model
- TestCleanRel
- mountTable
- .exeKeys
- github.com/Virgula0/app-listener/tools/codecheck
- bpfLsmListed
- WithSuspectsFrom
- github.com/google/fscrypt/actions.Context
- loadDaemonConfig
- SecretMatch
- Troubleshooting
- prctl_gate
- GenerateConf
- Vault
- IntegrationSuite
- taint_victim.c
- TestMain
- .requireSecretDenial
- stripImmutableFlags
- bootstrapSelfProtectPlaceholders
- partitionEncryptionRoots
- monitor: see what touches a path
- bleedLayout
- plan_test.go

## God Nodes (most connected - your core abstractions)
1. `Config` - 146 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 112 edges
4. `shell` - 97 edges
5. `Load()` - 95 edges
6. `absPath()` - 93 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 74 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `eBPF at three levels` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Why not an AppArmor profile?` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (252 total, 50 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.19
Nodes (24): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+16 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (167): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+159 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (11): Hash(), configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec() (+3 more)

### Community 3 - "op"
Cohesion: 0.15
Nodes (14): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+6 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (83): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+75 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (16): addErrKind, TestWithChmodDropWriteNeedsReadOnly(), classifyAddErr(), guardModeKey(), modeString(), addRecorder(), copySelf(), guardUnitTest (+8 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (34): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+26 more)

### Community 10 - "parseGuardEvent"
Cohesion: 0.11
Nodes (16): bpfGuardEvent, BusClient, AttestBusClient(), busClientKey(), execStamp(), imageOf(), pidfdPid(), procIdentity() (+8 more)

### Community 11 - "runWizard"
Cohesion: 0.05
Nodes (69): askBunTmpdirUsers(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), offerDaemonStart(), runDiffCatalog(), installFacultatives(), installSSHAgentFacultative() (+61 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (27): backendLink, configView, forwardForm, forwardView, installUI, protectedPage, recordReview, recordView (+19 more)

### Community 13 - "batchInodeRows"
Cohesion: 0.22
Nodes (12): noteReloadedSettings(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches() (+4 more)

### Community 14 - "Network"
Cohesion: 0.08
Nodes (11): binaryList, netGuardPill(), NetworkGuardRepository, newExeResolver(), guardRow(), clonePolicy(), NetExtra, NetPolicy (+3 more)

### Community 15 - "BinaryEntry"
Cohesion: 0.05
Nodes (31): deferredBinary, keptEntries, supersededBinary, BinaryEntry, checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath(), canonicalPaths() (+23 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (61): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+53 more)

### Community 18 - "fyne.io/fyne/v2/widget.Label"
Cohesion: 0.03
Nodes (23): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel, panelRenderer (+15 more)

### Community 19 - "fileEditModel"
Cohesion: 0.07
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (43): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+35 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (12): backingDeviceUnion(), Resource, BackingDevices(), GuardRepository, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, runtimeReleased() (+4 more)

### Community 23 - "auth_test.go"
Cohesion: 0.13
Nodes (38): CheckPassword(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf(), parseHash() (+30 more)

### Community 24 - "FileTarget"
Cohesion: 0.10
Nodes (12): filesPage, guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec() (+4 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.08
Nodes (79): add_inode_to_guard(), admit_port(), btrfs_copy_gate(), chmod_only_drops_write(), code_suspect(), count_degrade(), discover_guarded_parent(), emit_process_denial() (+71 more)

### Community 26 - "LibraryClosure"
Cohesion: 0.06
Nodes (33): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), IsScript(), ldConfInclude(), ldSoConfDirs() (+25 more)

### Community 27 - "go_pkg_os"
Cohesion: 0.08
Nodes (45): installGUI(), asRoot(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), ensureInstalledBinary(), ensureMasterKey(), fileSHA256(), prepareInstallation() (+37 more)

### Community 28 - "editView"
Cohesion: 0.12
Nodes (6): codeEntry, editView, fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 29 - "go_pkg_strings"
Cohesion: 0.07
Nodes (8): parseReservedPorts(), readReservedPorts(), TestParseReservedPorts(), Progress, User, fileID, i3Result, i3Workspace

### Community 30 - "RefuseIdentityLines"
Cohesion: 0.07
Nodes (26): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), systemPatterns(), systemRule(), RefuseIdentityLines(), whitelistLines(), ParseEventType() (+18 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 33 - "Guard"
Cohesion: 0.08
Nodes (11): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), rootHandle, openRoot(), SharedPinDegraded(), TestNoteRootMountFollowsMountID(), Guard (+3 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (21): keyEntry, listOpts, pickEntry, picker, treeSearch, lstatOf(), outsideHardlinks(), errText() (+13 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.26
Nodes (23): busPath(), connexec(), connFD(), createItem(), dial(), execInto(), failf(), get() (+15 more)

### Community 37 - "Message"
Cohesion: 0.08
Nodes (18): checkSecret(), entryOf(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp (+10 more)

### Community 38 - "controlServer"
Cohesion: 0.12
Nodes (12): authSelfPeer(), newControlManager(), readClientLines(), startControlManager(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout() (+4 more)

### Community 39 - "os.FileMode"
Cohesion: 0.11
Nodes (25): dupFD(), installFileAs(), installSSHAgent(), installSSHAgentEnv(), installUnit(), readAtNoFollow(), TestInstallSSHAgentRefusesSymlinkedUnitDir(), upsertFile() (+17 more)

### Community 40 - "RefreshSection"
Cohesion: 0.07
Nodes (35): logRefreshChange(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange(), applyLiveRefresh() (+27 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.05
Nodes (10): IntegrationSuite, IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, IntegrationSuite, launchCase (+2 more)

### Community 44 - "cursor.go"
Cohesion: 0.08
Nodes (21): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+13 more)

### Community 45 - "Candidate"
Cohesion: 0.11
Nodes (31): appendSectionsAndEdit(), collectDiffAdditions(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates() (+23 more)

### Community 48 - "update_test.go"
Cohesion: 0.09
Nodes (20): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+12 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (29): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+21 more)

### Community 51 - "Binary"
Cohesion: 0.12
Nodes (25): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), classifyLibs(), finishBinary() (+17 more)

### Community 52 - "Client"
Cohesion: 0.10
Nodes (11): Client, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines(), Reader (+3 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.08
Nodes (67): sameConfig(), runtimePathsFor(), auditGroupCoverage(), findResource(), groupWatchPaths(), watchGroup, addResource(), applyAllowLib() (+59 more)

### Community 54 - "session"
Cohesion: 0.10
Nodes (9): peerCred(), TestReadFrame(), decodeFrame(), encodeWithSerial(), header(), readFrame(), queued, senderQueue (+1 more)

### Community 55 - ".apply"
Cohesion: 0.14
Nodes (10): TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), applier, ledgerOps, PlannedBinary (+2 more)

### Community 56 - "githubRelease"
Cohesion: 0.21
Nodes (19): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), tags(), TestReleasesSince(), compareStableVersions(), filterChannel() (+11 more)

### Community 57 - "mutedText"
Cohesion: 0.08
Nodes (44): checkGrid, homePage, multiAnswer, recordForm, newLogo(), newBinaryList(), newBinaryRow(), newCheck() (+36 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (55): engine, portOwner, resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), nearestRootClaims() (+47 more)

### Community 59 - "Monitor"
Cohesion: 0.17
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 60 - "Prereq"
Cohesion: 0.20
Nodes (6): classifySetupError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), Prereq, fakeVault

### Community 61 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 62 - "inode_dev"
Cohesion: 0.15
Nodes (18): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), check_and_emit_args(), current_is_inspector(), exe_tag_bits() (+10 more)

### Community 64 - "policyEditor"
Cohesion: 0.13
Nodes (11): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+3 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "wm_test.go"
Cohesion: 0.10
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "dialLiveSession"
Cohesion: 0.14
Nodes (17): TestValidatePassword(), ValidatePassword(), dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig() (+9 more)

### Community 69 - "runOfflineEdit"
Cohesion: 0.10
Nodes (22): auditAfterEdit(), auditFindings(), auditFindingsWithConfig(), reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected() (+14 more)

### Community 70 - "wipe"
Cohesion: 0.17
Nodes (12): isFileVaultCiphertext(), isRegularFileTarget(), hasEncryptionPolicy(), isLockedRegularFileErr(), isNotEncryptedErr(), newBoundedKeyFn(), readKey(), TestIsNotEncryptedErr() (+4 more)

### Community 71 - "IntegrationSuite"
Cohesion: 0.16
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 72 - "journal"
Cohesion: 0.13
Nodes (7): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal, NetGuardConfig

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 75 - "LaunchRule"
Cohesion: 0.12
Nodes (20): allowedEnv(), binaryLaunchRules(), launchRulesByPath(), TestAllowedEnvKeepsOnlyRootPlaced(), TestLaunchRulesByPath(), argvEntry, applyLaunchRule(), envName() (+12 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "NetworkMonitorRepository"
Cohesion: 0.18
Nodes (4): NetworkMonitorRepository, NetworkMonitorUseCase, NewNetworkMonitorUseCase(), StreamUseCase

### Community 78 - "User"
Cohesion: 0.09
Nodes (23): entryWriters(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), bunLaunchersOffered() (+15 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "os.File"
Cohesion: 0.03
Nodes (87): multicallAdmissible(), ownerPaths(), scriptIdentity(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir() (+79 more)

### Community 81 - "recordFlow"
Cohesion: 0.13
Nodes (17): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+9 more)

### Community 82 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 83 - "rcUser"
Cohesion: 0.09
Nodes (54): listFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH(), removeFacultative(), removeFacultatives(), TestRemoveSSHAgentFacultativeKeepsGuardedSSHConfig() (+46 more)

### Community 84 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.13
Nodes (20): placedWalk, treeJudge, mountinfoDev(), SuperblockDev(), TestMountinfoDev(), confineBelow(), inUserHome(), vouchedFd() (+12 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.21
Nodes (16): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+8 more)

### Community 87 - "loginServer"
Cohesion: 0.11
Nodes (16): markReady(), markReadyWithLogin(), notifySystemdReady(), hasSentinel(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel(), startLoginServer() (+8 more)

### Community 88 - "vaultfs_test.go"
Cohesion: 0.12
Nodes (22): DefaultNewFileMode(), FileModeToUnixMode(), IsBinaryContent(), IsOctalMode(), ParseMode(), inodeOf(), mustRead(), mustWrite() (+14 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runView"
Cohesion: 0.07
Nodes (10): daemonPage, logGrid, runHost, runView, journalLine(), TestJournalLine(), newLogGrid(), captureLogs() (+2 more)

### Community 92 - "downloadAndVerify"
Cohesion: 0.33
Nodes (10): installHooks(), staged(), downloadAndVerify(), FetchLatest(), fetchPair(), FetchVerified(), Release, sanityCheckBinary() (+2 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "NetGuard"
Cohesion: 0.10
Nodes (11): PumpRingbuf(), reasonText(), eventTypeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), TestReasonText() (+3 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "time.Duration"
Cohesion: 0.15
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 100 - "hashFileShared"
Cohesion: 0.19
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 101 - "root.go"
Cohesion: 0.07
Nodes (28): releaseBuild(), Run(), closeDumpFile(), Execute(), initLogger(), postRun(), stageItem, stageState (+20 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "hashCache"
Cohesion: 0.13
Nodes (17): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, BinaryStat, DurableStat, HashCache (+9 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "newTestVetter"
Cohesion: 0.10
Nodes (28): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+20 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "highlight_test.go"
Cohesion: 0.10
Nodes (33): newChangelogModel(), rel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog(), TestNewChangelogModelViewContainsTitleAndNotes() (+25 more)

### Community 111 - "ConfinedEntry"
Cohesion: 0.29
Nodes (18): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+10 more)

### Community 112 - "installPage"
Cohesion: 0.07
Nodes (18): anyDaemonRunning(), actionGate, binSource, handButton, handCheck, installAction, installedBinary, installPage (+10 more)

### Community 113 - "Vault"
Cohesion: 0.24
Nodes (3): Vault, Open(), Unopened()

### Community 115 - "ListUsers"
Cohesion: 0.07
Nodes (34): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), grantNewBinaries() (+26 more)

### Community 116 - "deploy"
Cohesion: 0.20
Nodes (15): archAssets(), assetsFor(), checkUpdatePreconditions(), deploy(), fetchGUI(), guiAssetsFor(), hasAsset(), resolveAssets() (+7 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "UI"
Cohesion: 0.09
Nodes (22): showChangelog(), InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, InstallHooks (+14 more)

### Community 121 - "Recorder"
Cohesion: 0.16
Nodes (15): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start() (+7 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "inet_connect.c"
Cohesion: 0.26
Nodes (15): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+7 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "classify"
Cohesion: 0.22
Nodes (12): call(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyOpenCallsExact(), arg(), attributesOf(), classify(), classifyCollection() (+4 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "catalogRefresher"
Cohesion: 0.05
Nodes (34): drain(), newCatalogRefresher(), readEvents(), swapFile(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink() (+26 more)

### Community 131 - "audit.go"
Cohesion: 0.14
Nodes (17): auditEntry(), auditTree(), coveredByWatch(), relOrPath(), TestAuditTree(), TestWithin(), within(), readPutContent() (+9 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "changelog.go"
Cohesion: 0.06
Nodes (39): findButton, lineStyles, viewportState, findAll(), lowerRunes(), NewFinder(), newFindInput(), overlay() (+31 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "planUpdaters"
Cohesion: 0.14
Nodes (18): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+10 more)

### Community 143 - "serve.go"
Cohesion: 0.09
Nodes (30): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+22 more)

### Community 144 - "liveSession"
Cohesion: 0.16
Nodes (13): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck() (+5 more)

### Community 145 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.08
Nodes (30): class, pos, span, diffModel, editorModel, classOf(), column(), dedent() (+22 more)

### Community 146 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 148 - "GuardEvent"
Cohesion: 0.04
Nodes (27): holdHeadless(), HeadlessLine(), snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), commMatchesGuardedBinary(), GuardEvent, logBacklogDenial() (+19 more)

### Community 149 - "buildTrustedSet"
Cohesion: 0.09
Nodes (25): reserveChain(), inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys() (+17 more)

### Community 150 - "os.FileInfo"
Cohesion: 0.17
Nodes (6): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow()

### Community 152 - "bufio.Reader"
Cohesion: 0.11
Nodes (19): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+11 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.12
Nodes (26): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards(), selfGuards (+18 more)

### Community 155 - "dialog_test.go"
Cohesion: 0.36
Nodes (12): pickMode, openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder() (+4 more)

### Community 156 - "copytree.go"
Cohesion: 0.41
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 157 - "confirmUpdate"
Cohesion: 0.18
Nodes (11): channelOf(), confirmUpdate(), isTerminal(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), TestRunWithRefusesConcurrentRun() (+3 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 161 - "RootPlacedName"
Cohesion: 0.24
Nodes (7): RootPlacedName(), rootOnlyTree(), TestRootPlacedName(), userAssignable(), linkTarget(), entryKind, linkView

### Community 162 - "NetEventType"
Cohesion: 0.10
Nodes (14): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), FormatAddr(), NetEventTypes(), Ntohs() (+6 more)

### Community 163 - "secretJudge"
Cohesion: 0.12
Nodes (11): placedKey(), secretJudge, userGate, SecretService, LiveGuardAt(), Caller, Gate, Role (+3 more)

### Community 164 - "Row"
Cohesion: 0.29
Nodes (10): candidates(), confirmAll(), describe(), record(), run(), Given(), Row, Inspect() (+2 more)

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (45): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+37 more)

### Community 166 - "shell"
Cohesion: 0.04
Nodes (56): fakeLink, pageDef, pageID, reel, resizeWatch, shell, snapEngines, snapNet (+48 more)

### Community 167 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 168 - "procstats.go"
Cohesion: 0.38
Nodes (8): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic()

### Community 169 - "ensurePortAnchors"
Cohesion: 0.20
Nodes (10): ensureAnchor(), ensurePortAnchors(), ensureRootDir(), ephemeralRange(), ensureSecretAnchors(), PortAnchor(), TestParsePortSection(), isAbsoluteCandidate() (+2 more)

### Community 170 - "SelectOrphans"
Cohesion: 0.30
Nodes (13): appProtectorSet(), isAppProtector(), orphanedPolicies(), orphanedProtectors(), SelectOrphans(), appProtector(), foreignProtector(), policy() (+5 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "RunWith"
Cohesion: 0.07
Nodes (40): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), deletePostBackups(), listBackups(), restoreBackups() (+32 more)

### Community 175 - "TrustGuard"
Cohesion: 0.05
Nodes (22): exeRow, lsmHook, statFunc, tempGrant, tempMaskOp, tempRow, trustHook, trustRows (+14 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 178 - "run"
Cohesion: 0.36
Nodes (8): InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), syncGUI(), TestReadInstalledVersion(), GUIInstalled()

### Community 179 - "NetGuardEvent"
Cohesion: 0.06
Nodes (32): DrainUntilSignal(), MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runGuardHeadless(), runHeadless() (+24 more)

### Community 180 - "fileedit_symlink_test.go"
Cohesion: 0.36
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 186 - "Binary trust"
Cohesion: 0.22
Nodes (9): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch rules, Launch scan, Multicall binaries, The dynamic loader (+1 more)

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "github.com/spf13/cobra.Command"
Cohesion: 0.06
Nodes (39): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+31 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 195 - "Run"
Cohesion: 0.33
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "New"
Cohesion: 0.07
Nodes (52): decryptDirectories(), TestDecryptDirectoriesEmpty(), classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord() (+44 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "sync.Mutex"
Cohesion: 0.04
Nodes (85): CheckBPFLSM(), CheckEBPF(), newBinaryVetter(), openBinaryVetter(), buildGuards(), configureDaemonLogging(), confirmMasterKeyOverwrite(), lockOneRoot() (+77 more)

### Community 202 - "time.Time"
Cohesion: 0.07
Nodes (45): LockedError, ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI() (+37 more)

### Community 203 - "io.Writer"
Cohesion: 0.12
Nodes (11): formatElapsed(), orNone(), printBinary(), printRecording(), TestProgressReader(), fyneLogFilter, bottomBar, newBottomBar() (+3 more)

### Community 204 - "EventType"
Cohesion: 0.07
Nodes (19): selfEventsFor(), snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+11 more)

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (30): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+22 more)

### Community 209 - "tempGrantSetup"
Cohesion: 0.33
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 210 - "lockAndDeprovision"
Cohesion: 0.27
Nodes (7): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "KernelDev"
Cohesion: 0.16
Nodes (17): fdStat, physicalParent(), statFD(), KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices(), codeVouched() (+9 more)

### Community 213 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "askSSHAgentUsers"
Cohesion: 0.15
Nodes (16): addedResources(), TestAddedResourcesSkipsAlreadyGuardedSSH(), agentSocketNote(), askSSHAgentUsers(), isInsidePath(), CoveredPaths(), PathCovered(), SSHAgentSocketGuarded() (+8 more)

### Community 219 - "openLedger"
Cohesion: 0.29
Nodes (10): ConfirmGiven(), ConfirmPending(), refuses(), stillHolds(), openLedger(), TestConfirmGivenNeedsTheReviewedContent(), TestConfirmPendingNeverOpensALineTheLedgerDoesNotRefuse(), TestConfirmPendingOnlyWhatWasRefusedAndIsStillThere() (+2 more)

### Community 220 - "fakeVault"
Cohesion: 0.42
Nodes (3): TestDaemonUseCaseReloadRefusedDuringShutdown(), TestDaemonUseCaseStopBlocksUntilAllLocked(), fakeVault

### Community 221 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 223 - "record"
Cohesion: 0.22
Nodes (6): fixtureExe(), privateTempRoots(), record(), TestClassifyLibs(), TestClassifyRefusesWorldWritableLocations(), TestObserveKeysByKernelExeAndSkipsSelf()

### Community 224 - "ReadInstallCert"
Cohesion: 0.36
Nodes (8): InstallCert, ReadInstallCert(), RemoveInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip(), TestInstallCertSymlinkRefuses(), TestInstallCertUnreadableRefuses(), withCertPath()

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "github.com/charmbracelet/bubbletea.Model"
Cohesion: 0.04
Nodes (43): recordTickCmd(), runNetworkMonitor(), runTUI(), recordModel, formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision() (+35 more)

### Community 228 - "mountTable"
Cohesion: 0.18
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), userMountable() (+1 more)

### Community 231 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 232 - "WithSuspectsFrom"
Cohesion: 0.43
Nodes (3): WithSuspectsFrom(), netguardHarness, Option

### Community 233 - "github.com/google/fscrypt/actions.Context"
Cohesion: 0.40
Nodes (4): modifiedContextWithSource(), CleanOrphans(), listPolicies(), listProtectors()

### Community 234 - "loadDaemonConfig"
Cohesion: 0.40
Nodes (5): loadDaemonConfig(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 235 - "SecretMatch"
Cohesion: 0.40
Nodes (4): SecretScope, SecretMatch, parseSecretMatch(), TestSecretMatch()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 238 - "GenerateConf"
Cohesion: 0.05
Nodes (59): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+51 more)

### Community 242 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 244 - "stripImmutableFlags"
Cohesion: 0.50
Nodes (4): stripEntryFlags(), stripImmutableFlags(), TestStripImmutableFlagsDoesNotHangOnFIFO(), TestStripImmutableFlagsSkipsDanglingSymlink()

### Community 245 - "bootstrapSelfProtectPlaceholders"
Cohesion: 1.00
Nodes (3): ensurePinStateFilePlaceholder(), ensurePlaceholder(), bootstrapSelfProtectPlaceholders()

### Community 249 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 256 - "plan_test.go"
Cohesion: 0.11
Nodes (25): isWhiteout(), opaque(), pathByte(), pathRecord(), scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent() (+17 more)

## Knowledge Gaps
- **33 isolated node(s):** `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 664 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **50 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `TestMain`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `.requirePortDenied`, `go_pkg_strings`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.016787065703474373 - nodes in this community are weakly interconnected._
- **Why does `monitor()` connect `monitor` to `time.Duration`, `secret_probe/main.go`, `How it works`, `Installation`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `secretJudge`, `GuardInodeKey`, `time.Duration`, `hashFileShared`, `guardUnitTest`, `sync.Mutex`, `time.Time`, `fakeDaemon`, `TrustGuard`, `os.File`, `BinaryEntry`, `GuardEvent`, `daemonconfig.go`, `buildOneGuard`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `daemonconfig_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06511761331038439 - nodes in this community are weakly interconnected._