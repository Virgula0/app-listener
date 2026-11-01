# Graph Report - app-listener  (2026-11-01)

## Corpus Check
- 564 files · ~1,795,776 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7353 nodes · 33730 edges · 259 communities (213 shown, 46 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 3008 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `acb134b9`
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
- GuardEvent
- time.Time
- Network
- NetworkMonitorRepository
- buildOneGuard
- github.com/charmbracelet/bubbletea.Cmd
- IntegrationSuite
- guard_trust.bpf.c
- newPanel
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- recorder/plan.go
- guard.bpf.c
- LibraryClosure
- globBuilder
- SanitizeText
- NewGuard
- codeedit.go
- TestMonitorUseCaseLifecycle
- fileedit_model_test.go
- Prereq
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- controlServer
- github.com/charmbracelet/bubbles/textarea.Model
- watchSet
- daemon_multicall_test.go
- desktop_test.go
- monitorUnitTest
- Candidate
- find_test.go
- Guard
- stdio
- journal
- fadeScroll
- IntegrationSuite
- mutedText
- Client
- installPage
- session
- applier
- BinaryEntry
- icon
- engine
- audit.go
- RawTarget
- runView
- encodeLaunchSet
- FileTarget
- unix_connect.c
- context.Context
- IntegrationSuite
- Recorder
- RefreshSection
- Vault
- Config
- planUpdaters
- fcntl
- fileedit_model.go
- update_test.go
- codecheck/main.go
- recordFlow
- newTestVetter
- process_vm_readv.c
- os.File
- Binary
- .filterEvent
- rcUser
- OpenSystemPlaced
- bpfstats/main.go
- exe_supersede.h
- KernelDev
- gui/session.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- CandidateDir
- Vault
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- safeio.go
- NetEventType
- Well-known bypass classes
- NetGuardEvent
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- expandPlaceholders
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- .posAt
- Caller
- .showSized
- runOfflineEdit
- IntegrationSuite
- ConfinedEntry
- DeployInstalledBinary
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- RunCmd
- IntegrationSuite
- WebSocket /ws client
- inet_connect.c
- install.sh
- hashCache
- classify
- make test-integration (rootful Docker suite)
- RunWith
- trust_test.go
- listFacultatives
- trace-app-libs.sh
- IntegrationSuite
- Highlighter
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- daemon: permanent protection with encryption at rest
- serve.go
- liveSession
- .existingWhitelist
- Ledger
- FileEvent
- fakeGuardRepo
- dialLiveSession
- Daemon configuration
- fakeStore
- mc_tag.h
- IntegrationSuite
- Monitor
- catalogRefresher
- CopyTreeWithProgress
- secretJudge
- Installation
- BinaryStat
- graphify-refresh.sh
- layout.go
- wm_test.go
- newTestTarget
- absPath
- feed.go
- shell
- inodeChain
- plan_test.go
- vaultfs.go
- SelectOrphans
- confirmUpdate
- Limitations
- fetchReleases
- parseDaemonEvents
- sync.Mutex
- net_tester/main.go
- procstats.go
- githubRelease
- Run
- bufio.Reader
- IntegrationSuite
- tempGrantSetup
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- github.com/spf13/cobra.Command
- Self-updating apps and the catalog
- fileedit_symlink_test.go
- os.FileInfo
- changelog_test.go
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- New
- pathCache
- pinstate_test.go
- DaemonEvent
- fakeDaemon
- GuardInodeKey
- sources
- monitor
- IntegrationSuite
- .populateWatchInodes
- backend.go
- .eventTypeAllowed
- EventType
- WithBottomBar
- .apply
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- linkView
- runNetworkGuard
- downloadAndVerify
- IntegrationSuite
- run
- refuseOutsideHardlinks
- VouchedDevs
- app-listener documentation
- supersede_probe.c
- eventUnitTest
- bleedLayout
- mountTable
- IntegrationSuite
- github.com/Virgula0/app-listener/tools/codecheck
- inode_dev.h
- Plan
- TestReadFrame
- ValidateTargets
- CleanupStalePins
- Troubleshooting
- deploy
- GenerateConf
- SystemTrusted
- IntegrationSuite
- taint_victim.c
- User
- Vault
- go_pkg_strings
- guardLSMHooks
- TestMain
- ReadInstallCert
- monitor: see what touches a path
- TestCleanRel
- runMonitor
- .netIntegrityContainer
- network monitor: see what a program talks to
- ensurePortAnchors
- daemonconfig.go
- IntegrationSuite
- .runEditorHarness
- FindDaemonProcesses

## God Nodes (most connected - your core abstractions)
1. `Config` - 147 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 112 edges
4. `Load()` - 97 edges
5. `shell` - 97 edges
6. `absPath()` - 93 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 75 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Modes` --references--> `monitor()`  [INFERRED]
  docs/README.md → integrationtests/secret_probe/main.go
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `Workspaces` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `eBPF at three levels` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (259 total, 46 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.18
Nodes (34): tempJournalRow, argvMatch, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion (+26 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (156): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+148 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (8): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote(), watchPathsInConfig()

### Community 3 - "op"
Cohesion: 0.17
Nodes (13): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+5 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.06
Nodes (84): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+76 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.09
Nodes (51): prctl_gate(), guard_task_prctl(), admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader() (+43 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (14): addErrKind, classifyAddErr(), eventMask(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+6 more)

### Community 9 - "mkdirs"
Cohesion: 0.12
Nodes (40): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), newWatchPattern(), systemBinaryPatterns(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped() (+32 more)

### Community 10 - "GuardEvent"
Cohesion: 0.08
Nodes (17): holdHeadless(), HeadlessLine(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary() (+9 more)

### Community 11 - "time.Time"
Cohesion: 0.09
Nodes (11): LockedError, sweepID, sweepMemo, codeEntry, editView, Entry, selectable(), fileSize() (+3 more)

### Community 12 - "Network"
Cohesion: 0.10
Nodes (14): netGuardPill(), NetworkGuardRepository, newExeResolver(), guardRow(), stopPumped(), clonePolicy(), NetEngines, NetExtra (+6 more)

### Community 13 - "NetworkMonitorRepository"
Cohesion: 0.15
Nodes (6): runHeadless(), ProtocolString(), NetworkMonitorRepository, NetworkMonitorUseCase, NewNetworkMonitorUseCase(), StreamUseCase

### Community 14 - "buildOneGuard"
Cohesion: 0.06
Nodes (51): CheckEBPF(), backingDeviceUnion(), buildGuards(), buildOneGuard(), eventFilterOptions(), newPinGeneration(), prepareDaemonStart(), reloadOnce() (+43 more)

### Community 15 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (46): recordTickCmd(), runNetworkMonitor(), runTUI(), recordModel, diffModel, editorModel, formatGuardEventLine(), sanitizeTerminalText() (+38 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.08
Nodes (65): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+57 more)

### Community 18 - "newPanel"
Cohesion: 0.03
Nodes (27): appTheme, binaryRow, clickPanel, homePage, iconTile, iconTileRenderer, navItem, navItemRenderer (+19 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (13): newFindInput(), clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser (+5 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (51): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+43 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (9): Resource, GuardRepository, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots(), vaultOpForGuard() (+1 more)

### Community 23 - "auth.go"
Cohesion: 0.13
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (95): exe_refused(), inode_dev(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), bus_key(), check_and_emit_args(), chmod_only_drops_write() (+87 more)

### Community 26 - "LibraryClosure"
Cohesion: 0.09
Nodes (22): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths(), LibraryClosure() (+14 more)

### Community 27 - "globBuilder"
Cohesion: 0.09
Nodes (19): entryWriters(), reserveChain(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, GlobChild, statFunc, compileGlobNames() (+11 more)

### Community 28 - "SanitizeText"
Cohesion: 0.06
Nodes (26): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), scriptIdentity(), appRunning(), libDirSig(), formatElapsed() (+18 more)

### Community 29 - "NewGuard"
Cohesion: 0.15
Nodes (15): DrainUntilSignal(), displayEntries(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), TestWithChmodDropWriteNeedsReadOnly(), canonicalPaths() (+7 more)

### Community 30 - "codeedit.go"
Cohesion: 0.12
Nodes (35): class, classOf(), Edit(), leadingIndent(), New(), onlySpaces(), SetText(), softTabBack() (+27 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "fileedit_model_test.go"
Cohesion: 0.20
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 33 - "Prereq"
Cohesion: 0.18
Nodes (7): confirmRunPrereq(), classifySetupError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), Prereq, fakeVault

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.22
Nodes (25): busPath(), connexec(), connFD(), createItem(), dial(), execInto(), failf(), fakebus() (+17 more)

### Community 37 - "Message"
Cohesion: 0.10
Nodes (14): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+6 more)

### Community 38 - "controlServer"
Cohesion: 0.16
Nodes (10): authSelfPeer(), readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), configEditor, controlServer (+2 more)

### Community 39 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.14
Nodes (14): findButton, pos, column(), dedent(), lines(), moveTo(), Navigate(), repeatKey() (+6 more)

### Community 40 - "watchSet"
Cohesion: 0.19
Nodes (7): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, Process inspectors

### Community 41 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 42 - "desktop_test.go"
Cohesion: 0.04
Nodes (47): TestEveryInstallModeMapped(), buttonCursor, handCheckGroup, handRadioGroup, handSelect, multiAnswer, stageItem, stageState (+39 more)

### Community 44 - "Candidate"
Cohesion: 0.14
Nodes (22): appendSectionsAndEdit(), collectDiffAdditions(), runDiffCatalog(), TestUncoveredCandidates(), uncoveredCandidates(), addManualDirectories(), editConfig(), groupCandidates() (+14 more)

### Community 45 - "find_test.go"
Cohesion: 0.21
Nodes (18): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+10 more)

### Community 46 - "Guard"
Cohesion: 0.07
Nodes (14): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), portOwner, rootHandle, SectionScan, openRoot() (+6 more)

### Community 48 - "journal"
Cohesion: 0.15
Nodes (6): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 49 - "fadeScroll"
Cohesion: 0.08
Nodes (15): fadeBar, fadeBarRenderer, fadeScroll, fadeScrollRenderer, fadeScrollOf(), newTestFadeScroll(), settle(), TestFadeBarDragEndsGlide() (+7 more)

### Community 50 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 51 - "mutedText"
Cohesion: 0.05
Nodes (30): backendLink, checkGrid, configView, forwardForm, forwardView, protectedPage, recordForm, recordReview (+22 more)

### Community 52 - "Client"
Cohesion: 0.09
Nodes (14): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+6 more)

### Community 53 - "installPage"
Cohesion: 0.05
Nodes (40): anyDaemonRunning(), actionGate, binSource, handButton, handCheck, installAction, installedBinary, InstallMode (+32 more)

### Community 54 - "session"
Cohesion: 0.11
Nodes (9): userGate, Gate, Listen(), peerCred(), readFrame(), Bus, queued, senderQueue (+1 more)

### Community 55 - "applier"
Cohesion: 0.12
Nodes (9): Apply(), DeleteBackup(), Outcome, ProgressFunc, TestDeleteBackupOnlyItsOwn(), applier, ledgerOps, Store (+1 more)

### Community 56 - "BinaryEntry"
Cohesion: 0.07
Nodes (21): deferredBinary, keptEntries, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath() (+13 more)

### Community 57 - "icon"
Cohesion: 0.15
Nodes (6): binaryList, filesPage, newBinaryList(), newBinaryRow(), icon(), newButton()

### Community 58 - "engine"
Cohesion: 0.04
Nodes (29): engine, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+21 more)

### Community 59 - "audit.go"
Cohesion: 0.15
Nodes (19): auditEntry(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+11 more)

### Community 60 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 62 - "runView"
Cohesion: 0.06
Nodes (14): logGrid, logHook, logLine, runView, journalLine(), TestJournalLine(), logRow(), logStyle() (+6 more)

### Community 63 - "encodeLaunchSet"
Cohesion: 0.15
Nodes (15): allowedEnv(), TestAllowedEnvKeepsOnlyRootPlaced(), argvEntry, applyLaunchRule(), envName(), ParseAllowEnv(), parseLaunchRule(), BuiltinAllowedEnv() (+7 more)

### Community 64 - "FileTarget"
Cohesion: 0.07
Nodes (22): netGuardTab, policyEditor, policyHooks, policyTab, newFilesPage(), newPolicyTab(), verbFor(), netWarnings() (+14 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "context.Context"
Cohesion: 0.14
Nodes (12): detectWM(), centered(), FitSize(), Session, poll(), TestFitSize(), hyprClient, hyprMonitor (+4 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.21
Nodes (6): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent()

### Community 68 - "Recorder"
Cohesion: 0.14
Nodes (15): Target, Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start(), record() (+7 more)

### Community 69 - "RefreshSection"
Cohesion: 0.11
Nodes (24): logRefreshChange(), TestLogRefreshChange(), RefreshOptions, IsScript(), TestIsScript(), ScriptsOf(), TestScriptsOf_NamesTheConfiguredLine(), admitNew() (+16 more)

### Community 70 - "Vault"
Cohesion: 0.09
Nodes (29): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), encryptDirectories(), secureResources(), verifyEncryptionState(), resolveFilesystemPrereqs(), askEncryption() (+21 more)

### Community 71 - "Config"
Cohesion: 0.05
Nodes (52): ServeConfig, newControlManager(), startControlManager(), catchLifecycleSignals(), makeReloadHandler(), newDaemonModel(), runDaemon(), runDaemonUI() (+44 more)

### Community 72 - "planUpdaters"
Cohesion: 0.14
Nodes (18): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+10 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 75 - "update_test.go"
Cohesion: 0.12
Nodes (15): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFilterChannel(), testKeyPair(), TestNewerThanStable(), TestParseChecksum() (+7 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "recordFlow"
Cohesion: 0.15
Nodes (12): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+4 more)

### Community 78 - "newTestVetter"
Cohesion: 0.44
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.21
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "os.File"
Cohesion: 0.04
Nodes (76): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+68 more)

### Community 81 - "Binary"
Cohesion: 0.12
Nodes (25): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), classifyLibs(), finishBinary() (+17 more)

### Community 83 - "rcUser"
Cohesion: 0.12
Nodes (43): rcTarget, BunTmpDir(), TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv() (+35 more)

### Community 84 - "OpenSystemPlaced"
Cohesion: 0.07
Nodes (31): placedWalk, treeJudge, launchKey(), TestLaunchRulesFitTheKey(), FirstOnBtrfs(), OnBtrfs(), SuperblockDev(), TestSuperblockDevMatchesStatOffBtrfs() (+23 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.22
Nodes (14): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+6 more)

### Community 87 - "KernelDev"
Cohesion: 0.07
Nodes (29): markReadyWithLogin(), hasSentinel(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel(), runtimePathsFor(), startLoginServer(), loggedInUsers() (+21 more)

### Community 88 - "gui/session.go"
Cohesion: 0.20
Nodes (15): ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+7 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "CandidateDir"
Cohesion: 0.14
Nodes (18): bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), bunLaunchersOffered(), Discover(), DiscoverSystem(), CandidateDir, isAbsoluteCandidate() (+10 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "NetGuard"
Cohesion: 0.07
Nodes (20): foreignPinnedLink(), PumpRingbuf(), eventTypeKey(), mcNameKey(), statInodeKey(), codeVouched(), isTraced(), judgePasses() (+12 more)

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
Cohesion: 0.23
Nodes (9): duration(), formatClock(), New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse (+1 more)

### Community 100 - "safeio.go"
Cohesion: 0.08
Nodes (36): TestDaemonUnitWithMetadataOutput(), dupFD(), installFileAs(), installSSHAgent(), installSSHAgentEnv(), installUnit(), readAtNoFollow(), TestInstallSSHAgentRefusesSymlinkedUnitDir() (+28 more)

### Community 101 - "NetEventType"
Cohesion: 0.07
Nodes (17): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), BinariesSummary(), TestBinariesSummary(), FormatAddr() (+9 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "NetGuardEvent"
Cohesion: 0.15
Nodes (7): snapNetGuard, NetGuardEvent, appendCapped(), protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "expandPlaceholders"
Cohesion: 0.13
Nodes (8): symlinkTargetGlob(), expandPlaceholders(), TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - ".posAt"
Cohesion: 0.11
Nodes (12): lineStyles, viewportState, fg(), baseStyle(), lineLen(), markStyle(), repeatSpaces(), stylesOf() (+4 more)

### Community 111 - "Caller"
Cohesion: 0.18
Nodes (4): placedKey(), Caller, Role, Judge

### Community 112 - ".showSized"
Cohesion: 0.08
Nodes (19): showChangelog(), installUI, scrollable(), newFadeScroll(), ask(), diffText(), promptBody(), showAudit() (+11 more)

### Community 113 - "runOfflineEdit"
Cohesion: 0.10
Nodes (21): auditAfterEdit(), auditFindings(), reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive() (+13 more)

### Community 115 - "ConfinedEntry"
Cohesion: 0.10
Nodes (33): markReady(), notifySystemdReady(), binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared() (+25 more)

### Community 116 - "DeployInstalledBinary"
Cohesion: 0.07
Nodes (37): installGUI(), buildBinaryIfNeeded(), installBinaryOnly(), installConfig(), mustCwd(), removeSessionHooks(), applyUpdate(), deployGUI() (+29 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "RunCmd"
Cohesion: 0.12
Nodes (22): offerDaemonStart(), installFacultatives(), installSSHAgentFacultative(), runFacultatives(), TestInstallSSHAgentFacultative(), certifyInstall(), markInstallComplete(), preflightDeployedBinary() (+14 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "inet_connect.c"
Cohesion: 0.26
Nodes (15): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+7 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 127 - "classify"
Cohesion: 0.22
Nodes (12): call(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyOpenCallsExact(), arg(), attributesOf(), classify(), classifyCollection() (+4 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "RunWith"
Cohesion: 0.14
Nodes (23): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), deletePostBackups(), listBackups(), restoreBackups() (+15 more)

### Community 130 - "trust_test.go"
Cohesion: 0.19
Nodes (14): fakeLink, newFakeLink(), protectedShell(), TestForwardRefusesOwnExecutable(), TestProtectedEditFlow(), TestProtectedLiveCardsGated(), TestRecordReviewExtendsAnExistingSection(), confirmOpen() (+6 more)

### Community 131 - "listFacultatives"
Cohesion: 0.14
Nodes (18): addedResources(), TestAddedResourcesSkipsAlreadyGuardedSSH(), listFacultatives(), facultativeUser(), TestListFacultativesIgnoresDaemonConf(), TestListFacultativesMarksGuardedSSH(), isInsidePath(), CoveredPaths() (+10 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "Highlighter"
Cohesion: 0.14
Nodes (11): span, findAll(), lowerRunes(), runesEqual(), TestFindAllSmartCase(), Highlighter, lexerFor(), shebangLexer() (+3 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.13
Nodes (15): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+7 more)

### Community 143 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 144 - "liveSession"
Cohesion: 0.19
Nodes (11): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), eventsLabel(), forwardBinaries(), holdForward(), normalizedEvents() (+3 more)

### Community 145 - ".existingWhitelist"
Cohesion: 0.09
Nodes (25): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), confSafeMatch(), BinaryRule, homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob() (+17 more)

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (25): candidates(), confirmAll(), describe(), record(), run(), applySchema(), Ledger, Pending (+17 more)

### Community 147 - "FileEvent"
Cohesion: 0.12
Nodes (6): snapMonitor, FileEvent, MonitorRepository, newFakeMonitorRepo(), fakeSource, fakeMonitorRepo

### Community 148 - "fakeGuardRepo"
Cohesion: 0.05
Nodes (25): confirmUserPlaced(), forwardCheck(), binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), planTempBlock() (+17 more)

### Community 149 - "dialLiveSession"
Cohesion: 0.15
Nodes (16): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+8 more)

### Community 150 - "Daemon configuration"
Cohesion: 0.17
Nodes (12): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Launch rules: `[launch_rules]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Port sections, Reloading (+4 more)

### Community 152 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 154 - "Monitor"
Cohesion: 0.17
Nodes (5): NewMonitor(), putBtrfsLayout(), RealFileEngines(), Monitor, watchedInode

### Community 155 - "catalogRefresher"
Cohesion: 0.09
Nodes (19): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps() (+11 more)

### Community 156 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 157 - "secretJudge"
Cohesion: 0.15
Nodes (8): secretJudge, SecretScope, SecretMatch, SecretService, parseSecretMatch(), TestSecretMatch(), LiveGuardAt(), Scope

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "BinaryStat"
Cohesion: 0.21
Nodes (13): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+5 more)

### Community 161 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "wm_test.go"
Cohesion: 0.37
Nodes (13): reachable(), Detect(), hyprlandSocket(), i3RuntimeSocket(), socketOwnedBy(), envOf(), serve(), shortDir() (+5 more)

### Community 163 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 164 - "absPath"
Cohesion: 0.06
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (49): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+41 more)

### Community 166 - "shell"
Cohesion: 0.04
Nodes (45): daemonPage, pageDef, pageID, procStats, reel, resizeWatch, runHost, shell (+37 more)

### Community 167 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 168 - "plan_test.go"
Cohesion: 0.19
Nodes (19): isWhiteout(), opaque(), scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles() (+11 more)

### Community 169 - "vaultfs.go"
Cohesion: 0.12
Nodes (24): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), IsBinaryContent(), IsOctalMode(), ParseMode(), inodeOf(), mustRead() (+16 more)

### Community 170 - "SelectOrphans"
Cohesion: 0.20
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 171 - "confirmUpdate"
Cohesion: 0.18
Nodes (11): channelOf(), confirmUpdate(), isTerminal(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), TestRunWithRefusesConcurrentRun() (+3 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "fetchReleases"
Cohesion: 0.22
Nodes (8): downloadFile(), downloadReleaseFiles(), fetchReleases(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleases(), TestFetchReleasesHTTPError(), progressReader

### Community 174 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 175 - "sync.Mutex"
Cohesion: 0.07
Nodes (15): refusedReplacements, tempGrant, tempMaskOp, tempRow, trustRows, setSupersedeTrust(), liveTrust(), syncMap() (+7 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "procstats.go"
Cohesion: 0.30
Nodes (9): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+1 more)

### Community 178 - "githubRelease"
Cohesion: 0.21
Nodes (19): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+11 more)

### Community 179 - "Run"
Cohesion: 0.38
Nodes (7): planFrom(), newStage(), Run(), runSandbox(), tail(), Spec, stage

### Community 180 - "bufio.Reader"
Cohesion: 0.12
Nodes (18): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+10 more)

### Community 181 - "IntegrationSuite"
Cohesion: 0.23
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 182 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

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
Cohesion: 0.10
Nodes (31): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+23 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 194 - "os.FileInfo"
Cohesion: 0.21
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 195 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 199 - "New"
Cohesion: 0.07
Nodes (52): askDecrypt(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty(), classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey() (+44 more)

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "pinstate_test.go"
Cohesion: 0.09
Nodes (35): configureDaemonLogging(), loadDaemonConfig(), noteReloadedSettings(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive() (+27 more)

### Community 202 - "DaemonEvent"
Cohesion: 0.10
Nodes (32): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+24 more)

### Community 204 - "GuardInodeKey"
Cohesion: 0.06
Nodes (37): exeRow, inodeBatchFn, supersededBinary, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+29 more)

### Community 205 - "sources"
Cohesion: 0.21
Nodes (6): TestGeneratedLibraryBlockRoundTrips(), pathByte(), pathRecord(), TestPathRecord(), libDirDev(), sources

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (33): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+25 more)

### Community 210 - "backend.go"
Cohesion: 0.07
Nodes (51): asRoot(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), ensureInstalledBinary(), ensureMasterKey(), fileSHA256(), prepareInstallation(), runInstall() (+43 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "EventType"
Cohesion: 0.10
Nodes (18): selfEventsFor(), BusClient, EventType, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), AttestBusClient(), busClientKey(), execStamp() (+10 more)

### Community 213 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 214 - ".apply"
Cohesion: 0.29
Nodes (9): TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), TestPlanAdmitsOnlyWhatWasRecorded(), TestPlanExtendsAnExistingSection(), testRecording() (+1 more)

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "linkView"
Cohesion: 0.29
Nodes (6): linkTarget(), newLinkView(), TestCheckLinksJudgesEachTarget(), entryKind, libPlan, linkView

### Community 218 - "runNetworkGuard"
Cohesion: 0.11
Nodes (21): CheckBPFLSM(), confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), computeGuardBinaries(), confirmUnsafe(), discoverInfra() (+13 more)

### Community 219 - "downloadAndVerify"
Cohesion: 0.33
Nodes (10): installHooks(), staged(), downloadAndVerify(), FetchLatest(), fetchPair(), FetchVerified(), Release, sanityCheckBinary() (+2 more)

### Community 221 - "run"
Cohesion: 0.27
Nodes (10): fetchRelease(), InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), syncGUI(), TestFetchReleaseByTag() (+2 more)

### Community 222 - "refuseOutsideHardlinks"
Cohesion: 0.24
Nodes (8): lstatOf(), outsideHardlinks(), refuseOutsideHardlinks(), TestRefuseOutsideHardlinks_FileRoot(), TestRefuseOutsideHardlinks_InTreeLinksAccepted(), TestRefuseOutsideHardlinks_MissingRootFails(), TestRefuseOutsideHardlinks_OutsideLinkRefused(), writeFile()

### Community 223 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 224 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "eventUnitTest"
Cohesion: 0.27
Nodes (4): BpfEvent, eventUnitTest, Cstr(), DecodeBpfEvent()

### Community 228 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), userMountable() (+2 more)

### Community 231 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 232 - "Plan"
Cohesion: 0.38
Nodes (5): dedupe(), Plan, libOfChosen(), Extension, PlannedBinary

### Community 233 - "TestReadFrame"
Cohesion: 0.40
Nodes (4): TestReadFrame(), decodeFrame(), encodeWithSerial(), header()

### Community 234 - "ValidateTargets"
Cohesion: 0.15
Nodes (11): drain(), existingSection(), guardedMatch(), overlaps(), refusedPath(), sectionFor(), validateTarget(), ValidateTargets() (+3 more)

### Community 235 - "CleanupStalePins"
Cohesion: 0.40
Nodes (5): CleanupStalePins(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "deploy"
Cohesion: 0.20
Nodes (15): archAssets(), assetsFor(), checkUpdatePreconditions(), deploy(), fetchGUI(), guiAssetsFor(), hasAsset(), resolveAssets() (+7 more)

### Community 238 - "GenerateConf"
Cohesion: 0.05
Nodes (63): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+55 more)

### Community 239 - "SystemTrusted"
Cohesion: 0.29
Nodes (8): rootOwnedInode(), rootOwnedSafe(), SystemTrusted(), systemTool(), archiveMembers(), decompressed(), tarMembers(), TestArchiveMembersIgnoresNonArchives()

### Community 242 - "User"
Cohesion: 0.09
Nodes (43): askBunTmpdirUsers(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), cleanupBackups(), deploy(), deployFiles(), grantNewBinaries() (+35 more)

### Community 244 - "go_pkg_strings"
Cohesion: 0.05
Nodes (17): releaseBuild(), Run(), TestAuditTree(), TestWithin(), closeDumpFile(), Execute(), initLogger(), postRun() (+9 more)

### Community 245 - "guardLSMHooks"
Cohesion: 0.29
Nodes (5): lsmHook, trustHook, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), TestSocketHooksAreRequired()

### Community 246 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 247 - "ReadInstallCert"
Cohesion: 0.43
Nodes (7): InstallCert, ReadInstallCert(), RequireNoInstallCert(), TestInstallCertRoundTrip(), TestInstallCertSymlinkRefuses(), TestInstallCertUnreadableRefuses(), withCertPath()

### Community 251 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 254 - "runMonitor"
Cohesion: 0.20
Nodes (9): MakeDisplayPaths(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runHeadless(), runMonitor(), runTUI() (+1 more)

### Community 257 - "network monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 258 - "ensurePortAnchors"
Cohesion: 0.29
Nodes (8): ensureAnchor(), ensurePortAnchors(), ensureRootDir(), ephemeralRange(), parseReservedPorts(), readReservedPorts(), TestParseReservedPorts(), ensureSecretAnchors()

### Community 267 - "daemonconfig.go"
Cohesion: 0.08
Nodes (60): watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyBlockLine(), applyConfigLine(), applyDirective(), applyElectronApp() (+52 more)

### Community 294 - "FindDaemonProcesses"
Cohesion: 0.10
Nodes (21): applyLiveRefresh(), runRefreshCatalogOnly(), softenAutomatedRefreshErr(), refreshCatalogConfig(), TestApplyLiveRefreshDeliversReload(), TestSoftenAutomatedRefreshErr(), confirmUnit(), TestUnifiedDiff() (+13 more)

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 664 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **46 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `daemon_multicall_test.go`, `parseDaemonEvents`, `IntegrationSuite`, `IntegrationSuite`, `go_pkg_strings`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `TestMain`, `IntegrationSuite`, `.netIntegrityContainer`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.021776729559748426 - nodes in this community are weakly interconnected._
- **Why does `monitor()` connect `monitor` to `app-listener documentation`, `time.Duration`, `secret_probe/main.go`, `Installation`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0943859649122807 - nodes in this community are weakly interconnected._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `time.Duration`, `guardUnitTest`, `DaemonEvent`, `fakeDaemon`, `time.Time`, `GuardInodeKey`, `buildOneGuard`, `sync.Mutex`, `os.File`, `GuardEvent`, `ConfinedEntry`, `secretJudge`, `BinaryEntry`, `engine`, `NewGuard`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Should `daemonconfig_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06442577030812324 - nodes in this community are weakly interconnected._